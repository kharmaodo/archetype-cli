package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// runCommand exécute une commande shell et affiche la sortie en temps réel
func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// checkToolVersion vérifie si un outil est installé et affiche sa version
func checkToolVersion(tool string, versionArg string) (string, error) {
	cmd := exec.Command(tool, versionArg)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(out), "\n")
	if len(lines) > 0 {
		return strings.TrimSpace(lines[0]), nil
	}
	return "", nil
}

func main() {
	// Déclaration des flags
	install := flag.Bool("install", false, "Installer le jar de l'archetype localement")
	testGen := flag.Bool("test", false, "Tester la génération d'un projet Maven après installation")
	custom := flag.Bool("custom", false, "Personnaliser les valeurs Maven pour la génération")
	flag.Parse()

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("🔍 Vérification des outils requis...")

	javaVersion, err := checkToolVersion("java", "-version")
	if err != nil {
		fmt.Println("❌ Java n'est pas installé ou n'est pas accessible")
		os.Exit(1)
	}
	fmt.Println("✅ Java :", javaVersion)

	mavenVersion, err := checkToolVersion("mvn", "-v")
	if err != nil {
		fmt.Println("❌ Maven n'est pas installé ou n'est pas accessible")
		os.Exit(1)
	}
	fmt.Println("✅ Maven :", mavenVersion)

	// Vérification de l'existence du jar
	jarPath := "factory/usine.jar"
	if _, err := os.Stat(jarPath); os.IsNotExist(err) {
		fmt.Printf("❌ Fichier JAR introuvable : %s\n", jarPath)
		return
	}

	// Définition des valeurs par défaut
	groupId := "baobao"
	artifactId := "pi-zb"
	version := "1.0-SNAPSHOT"
	packageName := "sn.cbao"

	// Installation du jar si demandé
	if *install {
		fmt.Println("⚙️ Installation du jar en local...")
		err := runCommand("mvn", "install:install-file",
			"-Dfile="+jarPath,
			"-DgroupId=com.votreorganisation.archetypes",
			"-DartifactId=starter-kit-archetype",
			"-Dversion=0.0.1-SNAPSHOT",
			"-Dpackaging=jar",
		)
		if err != nil {
			fmt.Println("❌ Échec de l'installation :", err)
			os.Exit(1)
		}
		fmt.Println("✅ Installation réussie !")
		fmt.Println("\nVous pouvez maintenant exécuter : mvn archetype:generate -DarchetypeCatalog=local")
	}

	// Personnalisation des valeurs Maven si custom
	if *custom {
		fmt.Println("\n⚡ Personnalisation des valeurs Maven")

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

		fmt.Printf("Package (%s) : ", packageName)
		input, _ = reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input != "" {
			packageName = input
		}

		// Proposer la génération après saisie
		fmt.Print("\n💡 Voulez-vous générer un projet Maven avec ces paramètres maintenant ? (y/n) : ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)
		if strings.ToLower(choice) == "y" {
			fmt.Println("\n⚙️ Génération du projet Maven...")
			err := runCommand("mvn", "archetype:generate",
				"-DarchetypeCatalog=local",
				"-DarchetypeGroupId=com.votreorganisation.archetypes",
				"-DarchetypeArtifactId=starter-kit-archetype",
				"-DarchetypeVersion=0.0.1-SNAPSHOT",
				"-DgroupId="+groupId,
				"-DartifactId="+artifactId,
				"-Dversion="+version,
				"-Dpackage="+packageName,
				"-DinteractiveMode=false",
			)
			if err != nil {
				fmt.Println("❌ Échec de la génération :", err)
				os.Exit(1)
			}
			fmt.Println("✅ Projet Maven généré avec succès !")
		} else {
			fmt.Println("\nℹ️ Pour générer manuellement, utilisez la commande suivante :")
			fmt.Printf(`mvn archetype:generate \
  -DarchetypeCatalog=local \
  -DarchetypeGroupId=com.votreorganisation.archetypes \
  -DarchetypeArtifactId=starter-kit-archetype \
  -DarchetypeVersion=0.0.1-SNAPSHOT \
  -DgroupId=%s \
  -DartifactId=%s \
  -Dversion=%s \
  -Dpackage=%s \
  -DinteractiveMode=false
`, groupId, artifactId, version, packageName)
		}
	}

	// Test automatique si demandé avec --test
	if *testGen && !*custom {
		fmt.Println("\n⚙️ Test de génération du projet Maven avec les valeurs par défaut...")
		err := runCommand("mvn", "archetype:generate",
			"-DarchetypeCatalog=local",
			"-DarchetypeGroupId=com.votreorganisation.archetypes",
			"-DarchetypeArtifactId=starter-kit-archetype",
			"-DarchetypeVersion=0.0.1-SNAPSHOT",
			"-DgroupId="+groupId,
			"-DartifactId="+artifactId,
			"-Dversion="+version,
			"-Dpackage="+packageName,
			"-DinteractiveMode=false",
		)
		if err != nil {
			fmt.Println("❌ Échec de la génération :", err)
			os.Exit(1)
		}
		fmt.Println("✅ Projet Maven généré avec succès !")
	}
}
