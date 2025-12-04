package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func checkCommandExists(name string, versionArg string) (bool, string) {
	cmd := exec.Command(name, versionArg)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return false, ""
	}
	return true, strings.TrimSpace(string(out))
}

func checkProjectExists(dir string) bool {
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		return true
	}
	return false
}

func deleteProject(dir string) error {
	return os.RemoveAll(dir)
}

func installJar(jarPath string) error {
	cmd := exec.Command("mvn", "install:install-file",
		"-Dfile="+jarPath,
		"-DgroupId=com.votreorganisation.archetypes",
		"-DartifactId=starter-kit-archetype",
		"-Dversion=0.0.1-SNAPSHOT",
		"-Dpackaging=jar",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func generateProject(groupId, artifactId, version, pkg string) error {
	args := []string{
		"archetype:generate",
		"-DarchetypeCatalog=local",
		"-DarchetypeGroupId=com.votreorganisation.archetypes",
		"-DarchetypeArtifactId=starter-kit-archetype",
		"-DarchetypeVersion=0.0.1-SNAPSHOT",
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

func main() {
	installFlag := flag.Bool("install", false, "Installer le JAR de l'archetype")
	testFlag := flag.Bool("test", false, "Tester la génération après installation")
	customFlag := flag.Bool("custom", false, "Personnaliser groupId, artifactId, version et package")
	flag.Parse()

	jarPath := "factory/usine.jar"
	if _, err := os.Stat(jarPath); os.IsNotExist(err) {
		fmt.Printf("❌ Fichier JAR introuvable : %s\n", jarPath)
		return
	}

	fmt.Println("🔍 Vérification des outils requis...")
	javaOK, javaVersion := checkCommandExists("java", "-version")
	mvnOK, mvnVersion := checkCommandExists("mvn", "-v")

	if !javaOK || !mvnOK {
		fmt.Println("❌ Java et Maven doivent être installés pour continuer.")
		if !javaOK {
			fmt.Println("Java non trouvé")
		} else {
			fmt.Printf("Java trouvé: %s\n", javaVersion)
		}
		if !mvnOK {
			fmt.Println("Maven non trouvé")
		} else {
			fmt.Printf("Maven trouvé: %s\n", mvnVersion)
		}
		return
	}

	fmt.Printf("✅ Java: %s\n", javaVersion)
	fmt.Printf("✅ Maven: %s\n", mvnVersion)

	if *installFlag {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("⚙️ Installer le JAR de l'archetype localement ? (y/n) : ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "y" || input == "Y" {
			fmt.Println("⚙️ Installation du jar en local...")
			if err := installJar(jarPath); err != nil {
				fmt.Printf("❌ Erreur lors de l'installation du JAR: %v\n", err)
				return
			}
			fmt.Println("✅ Installation réussie !")
			fmt.Println("Vous pouvez maintenant exécuter : mvn archetype:generate -DarchetypeCatalog=local")
		} else {
			fmt.Println("⚠️ Installation annulée.")
		}
	}

	if *testFlag || *customFlag {
		reader := bufio.NewReader(os.Stdin)

		// Valeurs par défaut
		groupId := "baobao"
		artifactId := "pi-zb"
		version := "1.0-SNAPSHOT"
		pkg := "sn.cbao"

		if *customFlag {
			fmt.Println("⚡ Personnalisation des valeurs Maven (laisser vide pour défauts)")
			fmt.Printf("GroupId (%s) : ", groupId)
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			if input != "" {
				groupId = input
			}

			fmt.Printf("ArtifactId (%s) : ", artifactId)
			input, _ = reader.ReadString('\n')
			input = strings.TrimSpace(input)
			if input != "" {
				artifactId = input
			}

			fmt.Printf("Version (%s) : ", version)
			input, _ = reader.ReadString('\n')
			input = strings.TrimSpace(input)
			if input != "" {
				version = input
			}

			fmt.Printf("Package (%s) : ", pkg)
			input, _ = reader.ReadString('\n')
			input = strings.TrimSpace(input)
			if input != "" {
				pkg = input
			}
		}

		if checkProjectExists(artifactId) {
			fmt.Printf("⚠️ Le projet '%s' existe déjà.\n", artifactId)
			fmt.Print("Voulez-vous le supprimer et régénérer le projet ? (y/n) : ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			if input == "y" || input == "Y" {
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

		fmt.Println("⚡ Génération du projet Maven à partir de l'archetype...")
		if err := generateProject(groupId, artifactId, version, pkg); err != nil {
			fmt.Printf("❌ Erreur lors de la génération du projet : %v\n", err)
			return
		}
		fmt.Println("✅ Projet généré avec succès !")
		fmt.Printf("Vous pouvez maintenant ouvrir le projet '%s'\n", artifactId)
	}
}
