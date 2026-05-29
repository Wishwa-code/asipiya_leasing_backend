package database

import (
	"fmt"
	adminModels "garment-management-backend/internal/admin/models"
	"gorm.io/gorm"
)

func SeedDatabase(db *gorm.DB) error {
	// 1. Seed Vehicle Types, Makes, and Models if not already seeded
	var typeCount int64
	if err := db.Model(&adminModels.VehicleType{}).Count(&typeCount).Error; err != nil {
		return fmt.Errorf("failed to count vehicle types: %w", err)
	}

	if typeCount == 0 {
		// 1. Define vehicle types
		types := []adminModels.VehicleType{
			{VehicleTypeName: "Car", Description: "Sedans, Hatchbacks, Coupes, and Wagons", Status: "Active"},
			{VehicleTypeName: "Van", Description: "Passenger and Cargo Vans", Status: "Active"},
			{VehicleTypeName: "SUV / Jeep", Description: "Sports Utility Vehicles and Off-road Jeeps", Status: "Active"},
			{VehicleTypeName: "Cab / Crew Cab", Description: "Single and Double Cab Pickups", Status: "Active"},
			{VehicleTypeName: "Three-Wheeler", Description: "Tuk-Tuks and Auto-rickshaws", Status: "Active"},
			{VehicleTypeName: "Motorbike", Description: "Motorcycles, Scooters, and Mopeds", Status: "Active"},
			{VehicleTypeName: "Lorry / Truck", Description: "Light and Heavy Commercial Lorries/Trucks", Status: "Active"},
			{VehicleTypeName: "Bus", Description: "Passenger transport buses", Status: "Active"},
			{VehicleTypeName: "Tractor", Description: "Agricultural and industrial tractors", Status: "Active"},
		}

		typeMap := make(map[string]uint)
		for i := range types {
			if err := db.FirstOrCreate(&types[i], adminModels.VehicleType{VehicleTypeName: types[i].VehicleTypeName}).Error; err != nil {
				return fmt.Errorf("failed to seed vehicle type %s: %w", types[i].VehicleTypeName, err)
			}
			typeMap[types[i].VehicleTypeName] = types[i].ID
		}

		// 2. Define makes grouped by type
		makesSeed := []struct {
			MakeName string
			TypeName string
		}{
			// Cars
			{"Suzuki", "Car"},
			{"Toyota", "Car"},
			{"Honda", "Car"},
			{"Nissan", "Car"},
			{"Micro", "Car"},
			{"Hyundai", "Car"},
			{"Kia", "Car"},
			{"Mazda", "Car"},
			// Vans
			{"Toyota", "Van"},
			{"Nissan", "Van"},
			{"Suzuki", "Van"},
			{"Mazda", "Van"},
			// SUVs
			{"Toyota", "SUV / Jeep"},
			{"Honda", "SUV / Jeep"},
			{"Nissan", "SUV / Jeep"},
			{"Mitsubishi", "SUV / Jeep"},
			{"Land Rover", "SUV / Jeep"},
			// Cabs
			{"Toyota", "Cab / Crew Cab"},
			{"Mitsubishi", "Cab / Crew Cab"},
			{"Nissan", "Cab / Crew Cab"},
			{"Isuzu", "Cab / Crew Cab"},
			// Three-Wheelers
			{"Bajaj", "Three-Wheeler"},
			{"TVS", "Three-Wheeler"},
			{"Piaggio", "Three-Wheeler"},
			// Motorbikes
			{"Bajaj", "Motorbike"},
			{"TVS", "Motorbike"},
			{"Honda", "Motorbike"},
			{"Yamaha", "Motorbike"},
			{"Suzuki", "Motorbike"},
			// Lorries
			{"Tata", "Lorry / Truck"},
			{"Mahindra", "Lorry / Truck"},
			{"Isuzu", "Lorry / Truck"},
			{"Mitsubishi", "Lorry / Truck"},
			// Buses
			{"Tata", "Bus"},
			{"Ashok Leyland", "Bus"},
			{"Mitsubishi", "Bus"},
			{"Toyota", "Bus"},
			// Tractors
			{"Massey Ferguson", "Tractor"},
			{"Kubota", "Tractor"},
			{"Tafe", "Tractor"},
		}

		for _, m := range makesSeed {
			typeID, exists := typeMap[m.TypeName]
			if !exists {
				continue
			}
			makeRec := adminModels.VehicleMake{
				VehicleMake:   m.MakeName,
				VehicleTypeID: typeID,
				Status:        "Active",
			}
			if err := db.FirstOrCreate(&makeRec, adminModels.VehicleMake{
				VehicleMake:   m.MakeName,
				VehicleTypeID: typeID,
			}).Error; err != nil {
				return fmt.Errorf("failed to seed vehicle make %s for type %s: %w", m.MakeName, m.TypeName, err)
			}
		}

		// 3. Define models (linked to vehicle types, as per the schema model)
		modelsSeed := []struct {
			ModelName string
			TypeName  string
		}{
			// Cars
			{"Alto", "Car"},
			{"Wagon R", "Car"},
			{"Spacia", "Car"},
			{"Swift", "Car"},
			{"Celerio", "Car"},
			{"Vitz", "Car"},
			{"Aqua", "Car"},
			{"Axio", "Car"},
			{"Premio", "Car"},
			{"Allion", "Car"},
			{"Corolla", "Car"},
			{"Prius", "Car"},
			{"Yaris", "Car"},
			{"Fit", "Car"},
			{"Grace", "Car"},
			{"Civic", "Car"},
			{"Shuttle", "Car"},
			{"Insight", "Car"},
			{"Sunny", "Car"},
			{"Leaf", "Car"},
			{"Dayz", "Car"},
			{"Panda", "Car"},
			{"Elite", "Car"},
			{"Tucson", "Car"},
			{"Sportage", "Car"},
			// Vans
			{"HiAce", "Van"},
			{"TownAce", "Van"},
			{"LiteAce", "Van"},
			{"Noah", "Van"},
			{"Voxy", "Van"},
			{"Caravan", "Van"},
			{"Vanette", "Van"},
			{"Every", "Van"},
			{"Bongo", "Van"},
			// SUVs
			{"Land Cruiser", "SUV / Jeep"},
			{"Prado", "SUV / Jeep"},
			{"RAV4", "SUV / Jeep"},
			{"Rush", "SUV / Jeep"},
			{"C-HR", "SUV / Jeep"},
			{"Vezel", "SUV / Jeep"},
			{"CR-V", "SUV / Jeep"},
			{"X-Trail", "SUV / Jeep"},
			{"Montero", "SUV / Jeep"},
			{"Outlander", "SUV / Jeep"},
			{"Defender", "SUV / Jeep"},
			// Cabs
			{"Hilux", "Cab / Crew Cab"},
			{"L200", "Cab / Crew Cab"},
			{"Triton", "Cab / Crew Cab"},
			{"Navara", "Cab / Crew Cab"},
			{"D-Max", "Cab / Crew Cab"},
			// Three-Wheelers
			{"RE", "Three-Wheeler"},
			{"Maxima", "Three-Wheeler"},
			{"King", "Three-Wheeler"},
			{"Ape", "Three-Wheeler"},
			// Motorbikes
			{"Pulsar", "Motorbike"},
			{"Discover", "Motorbike"},
			{"Platina", "Motorbike"},
			{"CT100", "Motorbike"},
			{"Apache", "Motorbike"},
			{"Metro", "Motorbike"},
			{"Wego", "Motorbike"},
			{"Dio", "Motorbike"},
			{"Hornet", "Motorbike"},
			{"CD70", "Motorbike"},
			{"FZ", "Motorbike"},
			{"RayZR", "Motorbike"},
			{"Gixxer", "Motorbike"},
			// Lorries
			{"Ace (Dimo Batta)", "Lorry / Truck"},
			{"Super Ace", "Lorry / Truck"},
			{"Bolero", "Lorry / Truck"},
			{"Maxximo", "Lorry / Truck"},
			{"Elf", "Lorry / Truck"},
			{"Forward", "Lorry / Truck"},
			{"Canter", "Lorry / Truck"},
			// Buses
			{"LPT", "Bus"},
			{"LP", "Bus"},
			{"Viking", "Bus"},
			{"Fuso", "Bus"},
			{"Coaster", "Bus"},
			// Tractors
			{"Massey Ferguson 240", "Tractor"},
			{"Kubota L4508", "Tractor"},
			{"Tafe 45 DI", "Tractor"},
		}

		for _, m := range modelsSeed {
			typeID, exists := typeMap[m.TypeName]
			if !exists {
				continue
			}
			modelRec := adminModels.VehicleModel{
				VehicleModelName: m.ModelName,
				VehicleTypeID:    typeID,
				Status:           "Active",
			}
			if err := db.FirstOrCreate(&modelRec, adminModels.VehicleModel{
				VehicleModelName: m.ModelName,
				VehicleTypeID:    typeID,
			}).Error; err != nil {
				return fmt.Errorf("failed to seed vehicle model %s for type %s: %w", m.ModelName, m.TypeName, err)
			}
		}
		fmt.Println("✅ Database seeded with default Sri Lankan Vehicle Types, Makes, and Models")
	}

	// 4. Seed colors if not already seeded
	var colorCount int64
	if err := db.Model(&adminModels.Color{}).Count(&colorCount).Error; err != nil {
		return fmt.Errorf("failed to count colors: %w", err)
	}
	if colorCount == 0 {
		colors := []adminModels.Color{
			{ColorName: "Red", Status: "Active"},
			{ColorName: "Blue", Status: "Active"},
			{ColorName: "Green", Status: "Active"},
			{ColorName: "White", Status: "Active"},
			{ColorName: "Black", Status: "Active"},
			{ColorName: "Silver", Status: "Active"},
			{ColorName: "Grey", Status: "Active"},
			{ColorName: "Gold", Status: "Active"},
			{ColorName: "Yellow", Status: "Active"},
			{ColorName: "Orange", Status: "Active"},
			{ColorName: "Brown", Status: "Active"},
			{ColorName: "Beige", Status: "Active"},
		}
		for i := range colors {
			if err := db.FirstOrCreate(&colors[i], adminModels.Color{ColorName: colors[i].ColorName}).Error; err != nil {
				return fmt.Errorf("failed to seed color %s: %w", colors[i].ColorName, err)
			}
		}
		fmt.Println("✅ Database seeded with default active colors")
	}

	return nil
}
