package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	fmt.Println("=== Maven Archetype Installer CLI ===")

	// 1. Vérifier le fichier JAR
	jarPath := filepath.Join("factory", "usine.jar")
	if _, err := os.Stat(jarPath); os.IsNotExist(err) {
		fmt.Printf("❌ Fichier JAR non trouvé : %s\n", jarPath)
		return
	}
	fmt.Printf("✔ Fichier JAR trouvé : %s\n", jarPath)

	// 2. Vérifier Java
	javaVersionCmd := exec.Command("java", "-version")
	javaOut, err := javaVersionCmd.CombinedOutput()
	if err != nil {
		fmt.Println("❌ Java n'est pas installé ou non trouvé dans le PATH.")
		return
	}
	fmt.Printf("✔ Java détecté :\n%s\n", string(javaOut))

	// 3. Vérifier Maven
	mvnVersionCmd := exec.Command("mvn", "-v")
	mvnOut, err := mvnVersionCmd.CombinedOutput()
	if err != nil {
		fmt.Println("❌ Maven n'est pas installé ou non trouvé dans le PATH.")
		return
	}
	fmt.Printf("✔ Maven détecté :\n%s\n", string(mvnOut))

	reader := bufio.NewReader(os.Stdin)

	// 4. Installer le JAR
	fmt.Print("Voulez-vous installer le JAR dans votre repository local Maven ? (y/n) : ")
	installChoice, _ := reader.ReadString('\n')
	installChoice = strings.TrimSpace(strings.ToLower(installChoice))

	if installChoice == "y" || installChoice == "yes" {
		installCmd := exec.Command("mvn", "install:install-file",
			"-Dfile="+jarPath,
			"-DgroupId=com.votreorganisation.archetypes",
			"-DartifactId=starter-kit-archetype",
			"-Dversion=0.0.1-SNAPSHOT",
			"-Dpackaging=jar")
		installCmd.Stdout = os.Stdout
		installCmd.Stderr = os.Stderr

		fmt.Println("\n🚀 Installation du JAR...")
		if err := installCmd.Run(); err != nil {
			fmt.Println("❌ Erreur pendant l'installation du JAR :", err)
			return
		}
		fmt.Println("✔ Installation terminée avec succès !")
	} else {
		fmt.Println("⚠ Installation annulée par l'utilisateur.")
	}

	// 5. Afficher aide pour génération
	fmt.Print("\nVoulez-vous voir l'aide pour générer un projet Maven depuis l'archetype ? (y/n) : ")
	helpChoice, _ := reader.ReadString('\n')
	helpChoice = strings.TrimSpace(strings.ToLower(helpChoice))

	if helpChoice == "y" || helpChoice == "yes" {
		fmt.Println("\n💡 Exemple de commande Maven pour générer un projet :")
		fmt.Println(`
mvn archetype:generate \
  -DarchetypeCatalog=local \
  -DarchetypeGroupId=com.votreorganisation.archetypes \
  -DarchetypeArtifactId=starter-kit-archetype \
  -DarchetypeVersion=0.0.1-SNAPSHOT \
  -DgroupId=<groupId> \
  -DartifactId=<artifactId> \
  -Dversion=<version> \
  -Dpackage=<package> \
  -DinteractiveMode=false
`)
		// 6. Demander si l'utilisateur veut générer
		fmt.Print("\nVoulez-vous générer un projet Maven maintenant ? (y/n) : ")
		genChoice, _ := reader.ReadString('\n')
		genChoice = strings.TrimSpace(strings.ToLower(genChoice))

		if genChoice == "y" || genChoice == "yes" {
			// 7. Demander les valeurs dynamiques
			fmt.Print("Entrez groupId (default: baobao) : ")
			groupId, _ := reader.ReadString('\n')
			groupId = strings.TrimSpace(groupId)
			if groupId == "" {
				groupId = "baobao"
			}

			fmt.Print("Entrez artifactId (default: pi-zb) : ")
			artifactId, _ := reader.ReadString('\n')
			artifactId = strings.TrimSpace(artifactId)
			if artifactId == "" {
				artifactId = "pi-zb"
			}

			fmt.Print("Entrez version (default: 1.0-SNAPSHOT) : ")
			version, _ := reader.ReadString('\n')
			version = strings.TrimSpace(version)
			if version == "" {
				version = "1.0-SNAPSHOT"
			}

			fmt.Print("Entrez package (default: sn.cbao) : ")
			packageName, _ := reader.ReadString('\n')
			packageName = strings.TrimSpace(packageName)
			if packageName == "" {
				packageName = "sn.cbao"
			}

			// 8. Construire la commande Maven
			cmdArgs := []string{
				"archetype:generate",
				"-DarchetypeCatalog=local",
				"-DarchetypeGroupId=com.votreorganisation.archetypes",
				"-DarchetypeArtifactId=starter-kit-archetype",
				"-DarchetypeVersion=0.0.1-SNAPSHOT",
				"-DgroupId=" + groupId,
				"-DartifactId=" + artifactId,
				"-Dversion=" + version,
				"-Dpackage=" + packageName,
				"-DinteractiveMode=false",
			}

			genCmd := exec.Command("mvn", cmdArgs...)
			genCmd.Stdout = os.Stdout
			genCmd.Stderr = os.Stderr

			fmt.Println("\n🚀 Génération du projet Maven...")
			if err := genCmd.Run(); err != nil {
				fmt.Println("❌ Erreur lors de la génération :", err)
				return
			}
			fmt.Println("✔ Projet Maven généré avec succès !")
		} else {
			fmt.Println("⚠ Génération du projet annulée par l'utilisateur.")
		}
	} else {
		fmt.Println("⚠ Aide pour la génération annulée.")
	}

	fmt.Println("\n✅ CLI terminée.")
}
