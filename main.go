package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// --- Structures de Configuration ---

// Config représente la structure complète du fichier config.json
type Config struct {
	Project struct {
		JarPath   string `json:"jar_path"`
		Copyright string `json:"copyright"`
	} `json:"project"`
	Archetype struct {
		GroupID    string `json:"group_id"`
		ArtifactID string `json:"artifact_id"`
		Version    string `json:"version"`
	} `json:"archetype"`
	Defaults struct {
		GroupID    string `json:"group_id"`
		ArtifactID string `json:"artifact_id"`
		Version    string `json:"version"`
		PackageName string `json:"package_name"`
	} `json:"defaults"`
}

// Variable globale pour stocker la configuration chargée
var cfg Config 

// --- Fonctions de Configuration ---

// loadConfig lit le fichier config.json et le décode.
func loadConfig(filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("impossible de lire le fichier de configuration '%s': %w", filepath, err)
	}
	
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return fmt.Errorf("impossible de parser le fichier de configuration JSON: %w", err)
	}
	
	// Vérification de base
	if cfg.Project.JarPath == "" {
		return fmt.Errorf("le chemin du JAR (project.jar_path) ne peut pas être vide dans la configuration")
	}
	
	return nil
}

// --- Fonctions Métier ---

// checkCommandExists vérifie si une commande existe et retourne sa première ligne de version.
func checkCommandExists(name string, versionArg string) (bool, string) {
	cmd := exec.Command(name, versionArg)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return false, ""
	}
	// Nettoie et extrait la première ligne.
	output := strings.TrimSpace(string(out))
	if output == "" {
		return true, "Version introuvable (commande silencieuse)"
	}
	return true, strings.Split(output, "\n")[0]
}

// checkProjectExists vérifie si un répertoire existe déjà.
func checkProjectExists(dir string) bool {
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		return true
	}
	return false
}

// deleteProject supprime récursivement un répertoire.
func deleteProject(dir string) error {
	return os.RemoveAll(dir)
}

// installJar exécute la commande mvn install:install-file en utilisant les valeurs de cfg.
func installJar() error {
	cmd := exec.Command("mvn", "install:install-file",
		"-Dfile="+cfg.Project.JarPath,
		"-DgroupId="+cfg.Archetype.GroupID,
		"-DartifactId="+cfg.Archetype.ArtifactID,
		"-Dversion="+cfg.Archetype.Version,
		"-Dpackaging=jar",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// generateProject exécute la commande mvn archetype:generate en utilisant les valeurs de cfg.
func generateProject(groupId, artifactId, version, pkg string) error {
	args := []string{
		"archetype:generate",
		"-DarchetypeCatalog=local",
		"-DarchetypeGroupId=" + cfg.Archetype.GroupID,
		"-DarchetypeArtifactId=" + cfg.Archetype.ArtifactID,
		"-DarchetypeVersion=" + cfg.Archetype.Version,
		"-DgroupId=" + groupId,
		"-DartifactId=" + artifactId,
		"-Dversion=" + version,
		"-Dpackage=" + pkg,
		"-DinteractiveMode=false",
	}
	cmd := exec.Command("mvn", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// --- Fonction Principale ---

func main() {
	// 0. Chargement de la configuration
	if err := loadConfig("config.json"); err != nil {
		fmt.Printf("❌ Erreur critique de configuration: %v\n", err)
		os.Exit(1)
	}

	// Définition des flags
	installFlag := flag.Bool("install", false, "⚙️ Installer le JAR de l'archetype")
	testFlag := flag.Bool("test", false, "🧪 Tester la génération après installation")
	customFlag := flag.Bool("custom", false, "✏️ Personnaliser groupId, artifactId, version et package")

	// Redéfinition de l'usage avec les données de copyright du fichier config
	flag.Usage = func() {
		fmt.Printf("📣 Copyright : %s\n", cfg.Project.Copyright)
		fmt.Printf("🛠️  Usage de %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	// 1. Vérification du fichier JAR (chemin tiré de la configuration)
	fmt.Printf("\n1. Vérification du fichier JAR : %s\n", cfg.Project.JarPath)
	if _, err := os.Stat(cfg.Project.JarPath); os.IsNotExist(err) {
		fmt.Printf("❌ Fichier JAR introuvable : %s\n", cfg.Project.JarPath)
		return
	}
	fmt.Println("✅ Fichier JAR trouvé.")

	// 2. Vérification des outils requis
	fmt.Println("\n2. Vérification des outils requis (Java et Maven)...")
	javaOK, javaVersion := checkCommandExists("java", "-version")
	mvnOK, mvnVersion := checkCommandExists("mvn", "-v")

	if !javaOK || !mvnOK {
		fmt.Println("\n--- ❌ PRÉREQUIS MANQUANTS ---")
		fmt.Println("Java et Maven doivent être installés et accessibles dans votre PATH pour continuer.")
		if !javaOK {
			fmt.Println("Java non trouvé.")
		} else {
			fmt.Printf("Java trouvé (info: %s)\n", javaVersion)
		}
		if !mvnOK {
			fmt.Println("Maven non trouvé.")
		} else {
			fmt.Printf("Maven trouvé (info: %s)\n", mvnVersion)
		}
		return
	}

	fmt.Printf("✅ Java OK (info: %s)\n", javaVersion)
	fmt.Printf("✅ Maven OK (info: %s)\n", mvnVersion)

	// 3. Installation de l'archétype (si --install est utilisé)
	if *installFlag {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("\n⚙️ Installer le JAR de l'archetype localement ? (y/n) : ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if strings.EqualFold(input, "y") { // Utiliser EqualFold pour accepter 'y' ou 'Y'
			fmt.Println("⚙️ Installation du jar en local...")
			if err := installJar(); err != nil {
				fmt.Printf("❌ Erreur lors de l'installation du JAR: %v\n", err)
				return
			}
			fmt.Println("✅ Installation réussie !")
			fmt.Println("\nProchaine étape: Exécutez 'mvn archetype:generate -DarchetypeCatalog=local'")
		} else {
			fmt.Println("⚠️ Installation annulée.")
		}
	}

	// 4. Génération de projet (si --test ou --custom est utilisé)
	if *testFlag || *customFlag {
		reader := bufio.NewReader(os.Stdin)

		// Valeurs par défaut tirées de la configuration
		groupId := cfg.Defaults.GroupID
		artifactId := cfg.Defaults.ArtifactID
		version := cfg.Defaults.Version
		pkg := cfg.Defaults.PackageName

		if *customFlag {
			fmt.Println("\n⚡ Personnalisation des valeurs Maven (laisser vide pour défauts)")
			// Raccourcissement de la logique de prompt
			prompt := func(name, currentVal string) string {
				fmt.Printf("%s (%s) : ", name, currentVal)
				input, _ := reader.ReadString('\n')
				input = strings.TrimSpace(input)
				if input != "" {
					return input
				}
				return currentVal
			}
			
			groupId = prompt("GroupId", groupId)
			artifactId = prompt("ArtifactId", artifactId)
			version = prompt("Version", version)
			pkg = prompt("Package", pkg)
		}

		// Vérification de l'existence du répertoire de sortie
		if checkProjectExists(artifactId) {
			fmt.Printf("\n⚠️ Le projet '%s' existe déjà.\n", artifactId)
			fmt.Print("Voulez-vous le supprimer et régénérer le projet ? (y/n) : ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			if strings.EqualFold(input, "y") {
				if err := deleteProject(artifactId); err != nil {
					fmt.Printf("❌ Impossible de supprimer le projet : %v\n", err)
					return
				}
				fmt.Println("✅ Projet supprimé avec succès.")
			} else {
				fmt.Println("❌ Génération annulée par l'utilisateur.")
				return
			}
		}

		fmt.Println("\n⚡ Génération du projet Maven à partir de l'archétype...")
		if err := generateProject(groupId, artifactId, version, pkg); err != nil {
			fmt.Printf("❌ Erreur lors de la génération du projet : %v\n", err)
			return
		}
		fmt.Println("✅ Projet généré avec succès !")
		fmt.Printf("Vous pouvez maintenant ouvrir le projet '%s'\n", artifactId)
	}
	
	if !*installFlag && !*testFlag && !*customFlag {
		fmt.Println("\nAucun flag d'action (--install, --test, ou --custom) n'a été spécifié. Utilisez -h pour l'aide.")
	}

	fmt.Println("\n--- Fin de l'exécution ---")
}