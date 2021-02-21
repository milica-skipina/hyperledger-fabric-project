package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// Damage describes car damages
type Damage struct {
	// ID          string  `json:"ID"`
	Description string  `json:"description"`
	Cost        float64 `json:"cost"`
}

// User describes user details (car owner, repairman, ...)
type User struct {
	ID       string  `json:"ID"`
	Name     string  `json:"name"`
	Lastname string  `json:"lastname"`
	Email    string  `json:"email"`
	Money    float64 `json:"money"`
}

// Asset describes basic details of what makes up a simple asset (car)
type Asset struct {
	ID             string   `json:"ID"`
	Brand          string   `json:"brand"`
	Model          string   `json:"model"`
	Year           int      `json:"year"`
	Color          string   `json:"color"`
	OwnerID        string   `json:"owner"`
	Damages        []Damage `json:"damages"`
	AppraisedValue float64  `json:"appraisedValue`
}

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	users := []User{
		{ID: "user1", Name: "Marko", Lastname: "Markovic", Email: "marko.markovic@email.com", Money: 10000.00},
		{ID: "user2", Name: "Jovan", Lastname: "Jovanovic", Email: "jovan.jovanovic@email.com", Money: 5000.00},
		{ID: "user3", Name: "Lazar", Lastname: "Lazarevic", Email: "lazar.lazarevic@email.com", Money: 3750.00},
	}
	assets := []Asset{
		{ID: "asset1", Brand: "fiat", Model: "500L", Year: 2018, Color: "black", OwnerID: "user1", Damages: []Damage{}, AppraisedValue: 7000.00},
		{ID: "asset2", Brand: "audi", Model: "A6", Year: 2016, Color: "blue", OwnerID: "user2", Damages: []Damage{}, AppraisedValue: 5000.00},
		{ID: "asset3", Brand: "bmw", Model: "500L", Year: 2017, Color: "red", OwnerID: "user2", Damages: []Damage{}, AppraisedValue: 12000.00},
		{ID: "asset4", Brand: "ford", Model: "500L", Year: 2013, Color: "gray", OwnerID: "user1", Damages: []Damage{}, AppraisedValue: 7350.00},
		{ID: "asset5", Brand: "toyota", Model: "500L", Year: 2017, Color: "black", OwnerID: "user1", Damages: []Damage{}, AppraisedValue: 4600.00},
		{ID: "asset6", Brand: "opel", Model: "astra", Year: 2018, Color: "black", OwnerID: "user3", Damages: []Damage{}, AppraisedValue: 6300.00},
	}

	for _, user := range users {
		userJSON, err := json.Marshal(user)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(user.ID, userJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	for _, asset := range assets {
		assetJSON, err := json.Marshal(asset)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(asset.ID, assetJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, id string, brand string, model string, year int, color string, owner string, appraisedValue float64) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", id)
	}

	asset := Asset{
		ID:             id,
		Brand:          brand,
		Model:          model,
		Year:           year,
		Color:          color,
		OwnerID:        owner,
		AppraisedValue: appraisedValue,
		Damages:        []Damage{},
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*Asset, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var asset Asset
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// ReadUser returns the user stored in the world state with given id.
func (s *SmartContract) ReadUser(ctx contractapi.TransactionContextInterface, id string) (*User, error) {
	userJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if userJSON == nil {
		return nil, fmt.Errorf("the user %s does not exist", id)
	}

	var user User
	err = json.Unmarshal(userJSON, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateAsset updates an existing asset in the world state with provided parameters.
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, id string, color string, size int, owner string, appraisedValue float64) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	// overwriting original asset with new asset
	asset := Asset{
		ID:             id,
		Color:          color,
		OwnerID:        owner,
		AppraisedValue: appraisedValue,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

// DeleteAsset deletes an given asset from the world state.
func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

// TransferAsset updates the owner field of asset with given id in world state.
func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, id string, newOwner string, withDamage bool) error {
	asset, err := s.ReadAsset(ctx, id)
	if err != nil {
		return fmt.Errorf("Car not found")
	}
	if asset.OwnerID == newOwner {
		return fmt.Errorf("New owner is same as current")
	}
	owner, err := s.ReadUser(ctx, asset.OwnerID)
	if err != nil {
		return fmt.Errorf("Owner not found")
	}
	newO, err := s.ReadUser(ctx, newOwner)
	if err != nil {
		return fmt.Errorf("New owner not found")
	}
	length := len(asset.Damages)
	totalPrice := asset.AppraisedValue
	if length == 0 {
		asset.OwnerID = newOwner
	} else if withDamage {
		totalDamage := 0.0
		for i := 0; i < length; i++ {
			totalDamage = totalDamage + asset.Damages[i].Cost
		}
		asset.OwnerID = newOwner
		totalPrice = totalPrice - totalDamage
	} else {
		return fmt.Errorf("Car has unrepaired damages")
	}
	if newO.Money < totalPrice {
		return fmt.Errorf("Customer doesn't have enough money on his account")
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}
	owner.Money = owner.Money + totalPrice
	newO.Money = newO.Money - totalPrice
	ownerJSON, err := json.Marshal(owner)
	if err != nil {
		return err
	}
	newOwnerJSON, err := json.Marshal(newO)
	if err != nil {
		return err
	}
	ctx.GetStub().PutState(owner.ID, ownerJSON)
	ctx.GetStub().PutState(newOwner, newOwnerJSON)
	return ctx.GetStub().PutState(id, assetJSON)
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("asset", "user")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*Asset
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Asset
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}
	return assets, nil
}

// GetAllUsers returns all users found in world state
func (s *SmartContract) GetAllUsers(ctx contractapi.TransactionContextInterface) ([]*User, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all users in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("user", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var users []*User
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var user User
		err = json.Unmarshal(queryResponse.Value, &user)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

// ChangeAssetColor updates color of asset with given ID
func (s *SmartContract) ChangeAssetColor(ctx contractapi.TransactionContextInterface, id string, color string) error {
	asset, err := s.ReadAsset(ctx, id)
	if err != nil {
		return fmt.Errorf("Car not found")
	}
	asset.Color = color
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

// CreateAssetDamage issues a new damage to the asset in the world state with given details.
func (s *SmartContract) CreateAssetDamage(ctx contractapi.TransactionContextInterface, id string, description string, cost float64) error {
	asset, err := s.ReadAsset(ctx, id)
	if err != nil {
		return fmt.Errorf("Car not found")
	}
	damage := Damage{
		Description: description,
		Cost:        cost,
	}
	asset.Damages = append(asset.Damages, damage)
	totalCost := 0.0
	for _, damage := range asset.Damages {
		totalCost = totalCost + damage.Cost
	}
	if totalCost > asset.AppraisedValue {
		return ctx.GetStub().DelState(id)
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(id, assetJSON)
}

// RepairDamages removes all damages from asset with given ID
func (s *SmartContract) RepairDamages(ctx contractapi.TransactionContextInterface, id string, mechanic string) error {
	asset, err := s.ReadAsset(ctx, id)
	if err != nil {
		return fmt.Errorf("Car not found")
	}
	owner, err := s.ReadUser(ctx, asset.OwnerID)
	if err != nil {
		return fmt.Errorf("Owner not found")
	}
	repairman, err := s.ReadUser(ctx, mechanic)
	if err != nil {
		return fmt.Errorf("Repairman not found")
	}

	totalCost := 0.0
	for _, damage := range asset.Damages {
		totalCost = totalCost + damage.Cost
	}

	if owner.Money < totalCost {
		return fmt.Errorf("Owner doesn't have enough money on his account")
	}

	owner.Money = owner.Money - totalCost
	repairman.Money = repairman.Money + totalCost
	ownerJSON, err := json.Marshal(owner)
	if err != nil {
		return err
	}
	repairmanJSON, err := json.Marshal(repairman)
	if err != nil {
		return err
	}

	asset.Damages = []Damage{}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}
	ctx.GetStub().PutState(owner.ID, ownerJSON)
	ctx.GetStub().PutState(mechanic, repairmanJSON)
	return ctx.GetStub().PutState(id, assetJSON)
}

// FindAssets returns all assets by color and owner
func (s *SmartContract) FindAssets(ctx contractapi.TransactionContextInterface, color string, owner string) ([]*Asset, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("asset", "user")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*Asset
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Asset
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		if (asset.Color == color || color == "") && (owner == "" || asset.OwnerID == owner) {
			assets = append(assets, &asset)
		}
	}
	return assets, nil
}
