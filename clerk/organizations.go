package clerk

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/organization"
	"github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found, using system environment variables")
	}

	clerk.SetKey(os.Getenv("CLERK_SECRET_KEY"))
}

func GetAllOrganizations() ([]string, error) {
	var allOrgIDs []string
	limit := int64(100)
	offset := int64(0)

	for {
		resp, err := organization.List(context.Background(), &organization.ListParams{
			ListParams: clerk.ListParams{
				Limit:  &limit,
				Offset: &offset,
			},
		})

		if err != nil {
			return nil, err
		}

		for _, org := range resp.Organizations {
			allOrgIDs = append(allOrgIDs, org.ID)
		}

		if int64(len(resp.Organizations)) < limit {
			break
		}
		offset += limit
	}

	return allOrgIDs, nil
}

func GetUserOrganizations(userId string) (*clerk.OrganizationMembershipList, error) {
	orgMemberships, err := user.ListOrganizationMemberships(context.Background(), userId, &user.ListOrganizationMembershipsParams{})

	if err != nil {
		log.Printf("Error getting organization memberships: %v", err)
	}

	return orgMemberships, nil
}

func GetUserOrganizationId(userId string) (string, error) {
	orgMemberships, err := GetUserOrganizations(userId)
	if err != nil {
		return "", err
	}
	if len(orgMemberships.OrganizationMemberships) == 0 {
		return "", errors.New("no organization memberships found")
	}
	return orgMemberships.OrganizationMemberships[0].Organization.ID, nil
}

func GetOrganizationPublicMetadata(organizationId string) (map[string]interface{}, error) {
	organization, err := organization.Get(context.Background(), organizationId)
	if err != nil {
		return nil, err
	}

	var metadata map[string]interface{}
	if err := json.Unmarshal(organization.PublicMetadata, &metadata); err != nil {
		return nil, err
	}

	return metadata, nil
}

func UpdateOrganizationPublicMetadata(organizationId string, metadata map[string]interface{}) error {
	jsonData, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	rawMessage := json.RawMessage(jsonData)
	_, err = organization.Update(context.Background(), organizationId, &organization.UpdateParams{
		PublicMetadata: &rawMessage,
	})

	return err
}
