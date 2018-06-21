package main

import (
	"fmt"
	"log"
)

func main() {
	err := run("bug1943")
	if err != nil {
		panic(err)
	}
}

func run(prefix string) error {
	client, err := buildAzureClient()
	if err != nil {
		return fmt.Errorf("Error building Azure Client: %+v", err)
	}

	name := fmt.Sprintf("%sClassicReg", prefix)
	storageAccountName := fmt.Sprintf("%sstor", prefix)
	resourceGroupName := fmt.Sprintf("%s-resources", prefix)
	location := "West Europe"

	err = client.createResourceGroup(resourceGroupName, location)
	if err != nil {
		return err
	}

	err = client.createStorageAccount(storageAccountName, resourceGroupName, location)
	if err != nil {
		return err
	}

	storageAccountId, err := client.getStorageAccount(storageAccountName, resourceGroupName)
	if err != nil {
		return err
	}
	log.Printf("Storage Account ID (from Storage Account): %q", *storageAccountId)

	err = client.createClassicContainerRegistry(name, *storageAccountId, resourceGroupName, location)
	if err != nil {
		return err
	}

	storageAccountIdFromReg, err := client.getStorageAccountIdFromContainerRegistry(name, resourceGroupName)
	if err != nil {
		return err
	}

	log.Printf("Storage Account ID (from Registry): %q", *storageAccountIdFromReg)

	defer client.deleteClassicContainerRegistry(name, resourceGroupName)
	defer client.deleteStorageAccount(name, resourceGroupName)
	defer client.deleteResourceGroup(resourceGroupName)

	return nil
}
