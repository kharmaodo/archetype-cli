---

## 1️⃣ Installer les dépendances Go

On utilisera le package [`flag`](https://pkg.go.dev/flag) pour gérer les options CLI et [`os/exec`](https://pkg.go.dev/os/exec) pour exécuter les commandes système.

```bash
go mod init archetype-cli
go get
```

---

## 2️⃣ Structure du projet

```
archetype-cli/
├── factory/                        # Contient starter-kit-archetype-0.0.1-SNAPSHOT.jar
├── cmd/
│   └── main.go                     # Entrée principale
├── go.mod
├── go.sum
└── README.md
```

---

## 3️⃣ Code complet de `main.go`

```go
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Chemin vers le JAR de l'archetype
const jarPath = "factory/starter-kit-archetype-0.0.1-SNAPSHOT.jar"

// Variables pour flags CLI
var (
	installFlag bool
	testFlag    bool
)

// Vérifie si une commande existe et retourne sa version
func checkCommand(cmdName string, versionArg string) (string, error) {
	cmd := exec.Command(cmdName, versionArg)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// Exécute une commande shell et affiche la sortie
func runCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Demande une confirmation à l'utilisateur (y/n)
func confirm(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s (y/n) : ", prompt)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
		if input == "y" || input == "yes" {
			return true
		} else if input == "n" || input == "no" {
			return false
		}
	}
}

func main() {
	// Définition des flags CLI
	flag.BoolVar(&installFlag, "install", false, "Installer le JAR de l'archetype Maven")
	flag.BoolVar(&testFlag, "test", false, "Tester la génération d'un projet Maven à partir de l'archetype")
	flag.Parse()

	// Vérifier Java
	javaVer, err := checkCommand("java", "-version")
	if err != nil {
		fmt.Println("Java n'est pas installé ou introuvable. Arrêt.")
		os.Exit(1)
	}
	fmt.Println("✅ Java installé :")
	fmt.Println(javaVer)

	// Vérifier Maven
	mvnVer, err := checkCommand("mvn", "-v")
	if err != nil {
		fmt.Println("Maven n'est pas installé ou introuvable. Arrêt.")
		os.Exit(1)
	}
	fmt.Println("✅ Maven installé :")
	fmt.Println(mvnVer)

	// Installation du JAR
	if installFlag || confirm("Voulez-vous installer le JAR de l'archetype Maven localement ?") {
		fmt.Println("\n🔧 Installation du JAR...")
		err := runCommand("mvn", "install:install-file",
			"-Dfile="+jarPath,
			"-DgroupId=com.votreorganisation.archetypes",
			"-DartifactId=starter-kit-archetype",
			"-Dversion=0.0.1-SNAPSHOT",
			"-Dpackaging=jar")
		if err != nil {
			fmt.Println("❌ Erreur lors de l'installation du JAR :", err)
			os.Exit(1)
		}
		fmt.Println("✅ JAR installé avec succès !")
	}

	// Test génération projet
	if testFlag || confirm("Voulez-vous tester la génération d'un projet Maven à partir de l'archetype ?") {
		fmt.Println("\n💡 Affichage de l'aide de génération :")
		runCommand("mvn", "archetype:generate", "-DarchetypeCatalog=local")

		if confirm("Voulez-vous générer un projet Maven automatiquement avec des valeurs dynamiques ?") {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("groupId (default baobao) : ")
			groupId, _ := reader.ReadString('\n')
			groupId = strings.TrimSpace(groupId)
			if groupId == "" {
				groupId = "baobao"
			}

			fmt.Print("artifactId (default pi-zb) : ")
			artifactId, _ := reader.ReadString('\n')
			artifactId = strings.TrimSpace(artifactId)
			if artifactId == "" {
				artifactId = "pi-zb"
			}

			fmt.Print("version (default 1.0-SNAPSHOT) : ")
			version, _ := reader.ReadString('\n')
			version = strings.TrimSpace(version)
			if version == "" {
				version = "1.0-SNAPSHOT"
			}

			fmt.Print("package (default sn.cbao) : ")
			pkg, _ := reader.ReadString('\n')
			pkg = strings.TrimSpace(pkg)
			if pkg == "" {
				pkg = "sn.cbao"
			}

			fmt.Println("\n🚀 Génération du projet Maven...")
			err := runCommand("mvn", "archetype:generate",
				"-DarchetypeCatalog=local",
				"-DarchetypeGroupId=com.votreorganisation.archetypes",
				"-DarchetypeArtifactId=starter-kit-archetype",
				"-DarchetypeVersion=0.0.1-SNAPSHOT",
				"-DgroupId="+groupId,
				"-DartifactId="+artifactId,
				"-Dversion="+version,
				"-Dpackage="+pkg,
				"-DinteractiveMode=false")
			if err != nil {
				fmt.Println("❌ Erreur lors de la génération :", err)
				os.Exit(1)
			}
			fmt.Println("✅ Projet Maven généré avec succès !")
		}
	}
}
```

---

## 4️⃣ Fonctionnalités de cette version

1. **Affichage lisible des versions exactes de Java et Maven**
2. **Flags CLI** : `--install` et `--test` pour automatiser sans interaction
3. **Mode interactif** si les flags ne sont pas utilisés
4. **Saisie dynamique des paramètres** Maven : `groupId`, `artifactId`, `version`, `package`
5. **Automatisation CI/CD** : possibilité de passer tous les paramètres via les flags pour scripts

---

## 5️⃣ Exemples d’utilisation

### 5.1 Mode interactif

```bash
./archetype-cli
```

### 5.2 Installation automatique

```bash
./archetype-cli --install
```

### 5.3 Tester génération automatiquement

```bash
./archetype-cli --install --test
```

---
