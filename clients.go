package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/containerregistry/mgmt/2017-10-01/containerregistry"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-02-01/resources"
	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2018-02-01/storage"
)

func (c *azureClient) createClassicContainerRegistry(name, storageAccountId, resourceGroupName, location string) error {
	log.Printf("Creating Classic Container Registry..")

	ctx := context.TODO()
	registry := containerregistry.Registry{
		Location: &location,
		Sku: &containerregistry.Sku{
			Name: containerregistry.Classic,
			Tier: containerregistry.SkuTierClassic,
		},
		RegistryProperties: &containerregistry.RegistryProperties{
			StorageAccount: &containerregistry.StorageAccountProperties{
				ID: &storageAccountId,
			},
		},
	}

	future, err := c.classicRegistriesClient.Create(ctx, resourceGroupName, name, registry)
	if err != nil {
		return err
	}

	err = future.WaitForCompletion(ctx, c.classicRegistriesClient.Client)
	if err != nil {
		return err
	}

	return nil
}

func (c *azureClient) getStorageAccountIdFromContainerRegistry(name, resourceGroup string) (*string, error) {
	log.Printf("Retrieving Storage Account ID from Classic Container Registry..")

	ctx := context.TODO()
	acc, err := c.classicRegistriesClient.Get(ctx, resourceGroup, name)
	if err != nil {
		return nil, err
	}


	return acc.StorageAccount.ID, nil
}

func (c *azureClient) deleteClassicContainerRegistry(name string, resourceGroup string) error {
	log.Printf("Deleting Classic Container Registry..")

	ctx := context.TODO()
	future, err := c.classicRegistriesClient.Delete(ctx, resourceGroup, name)
	if err != nil {
		return err
	}

	err = future.WaitForCompletion(ctx, c.classicRegistriesClient.Client)
	if err != nil {
		return err
	}

	return nil
}

func (c *azureClient) createStorageAccount(name, resourceGroupName, location string) error {
	log.Printf("Creating Storage Account..")
	ctx := context.TODO()
	service := storage.AccountCreateParameters{
		Location: &location,
		Kind: storage.Storage,
		Sku: &storage.Sku{
			Name: storage.StandardLRS,
			Tier: storage.Standard,
		},
	}

	future, err := c.storageAccountsClient.Create(ctx, resourceGroupName, name, service)
	if err != nil {
		return err
	}

	err = future.WaitForCompletion(ctx, c.storageAccountsClient.Client)
	if err != nil {
		return err
	}

	return nil
}

func (c *azureClient) getStorageAccount(name, resourceGroupName string) (*string, error) {
	log.Printf("Retrieving Storage Account..")

	ctx := context.TODO()
	account, err := c.storageAccountsClient.GetProperties(ctx, resourceGroupName, name)
	if err != nil {
		return nil, err
	}

	return account.ID, nil
}

func (c *azureClient) deleteStorageAccount(name, resourceGroupName string) error {
	ctx := context.TODO()
	log.Printf("Deleting Storage Account..")
	_, err := c.storageAccountsClient.Delete(ctx, resourceGroupName, name)
	if err != nil {
		return err
	}

	return nil
}

func (c *azureClient) createResourceGroup(name, location string) error {
	ctx := context.TODO()
	group := resources.Group{
		Location: &location,
	}

	log.Printf("Creating Resource Group..")
	_, err := c.resourceGroupsClient.CreateOrUpdate(ctx, name, group)
	if err != nil {
		return fmt.Errorf("Error creating Resource Group %q: %+v", name, err)
	}

	return nil
}

func (c *azureClient) deleteResourceGroup(name string) error {
	ctx := context.TODO()
	log.Printf("Deleting Resource Group..")
	future, err := c.resourceGroupsClient.Delete(ctx, name)
	if err != nil {
		return err
	}

	log.Printf("Waiting for deletion of Resource Group to complete..")
	err = future.WaitForCompletion(ctx, c.resourceGroupsClient.Client)
	if err != nil {
		return err
	}

	return nil
}