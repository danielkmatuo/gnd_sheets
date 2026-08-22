package backend

import "fmt"

func Run() error {
	root := findRootDir()
	if root == "" {
		return fmt.Errorf("couldnt find root dir")
	}

	fmt.Printf("root: %v\n", root)

    err := ensureChildDirsExist(root); 
	if err != nil {
        return err
    }

    err = ensureEssentialFilesExist(root); 
	if err != nil {
        return err
    }

    configureServer()

    return runServer(":8080")
}
