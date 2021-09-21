package dd

import (
	"strings"
	"testing"
)

const dataFilePath = "./device-detection-cxx/device-detection-data/51Degrees-LiteV4.1.hash"

func initResourceManager(
	properties string,
	t *testing.T) *ResourceManager {
	// Setup a resource manager
	manager := NewResourceManager()
	// Create config per performance
	config := NewConfigHash(Balanced)
	config.SetAllowUnmatched(false)
	err := InitManagerFromFile(
		manager,
		*config,
		properties,
		dataFilePath,
	)
	if err != nil {
		t.Errorf("Failed to initialize resource manager from file \"%s\"\n.",
			dataFilePath)
	}
	return manager
}

func freeResourceManager(
	manager *ResourceManager,
	t *testing.T) {
	// Free the resource manager
	err := manager.Free()
	if err != nil {
		t.Error("Failed to free resource manager.")
	}
}

// This tests available properties are returned correctly
func TestAvailableProperties(t *testing.T) {
	data := struct {
		properties        string
		expectedProperty  string
		expectedAvailable int
	}{"IsMobile", "IsMobile", 1}
	manager := initResourceManager(data.properties, t)
	// Create results hash
	results, err := NewResultsHash(manager, 1, 0)
	if err != nil {
		t.Error("Failed to create results hash object.")
	}

	available, err := results.AvailableProperties()
	if err != nil {
		t.Error("Failed to obtain available properties.")
	}

	if len(available) != data.expectedAvailable {
		t.Errorf("Expected to have %d properties but get %d.\n",
			data.expectedAvailable,
			len(available))
	}

	if !strings.EqualFold(data.expectedProperty, available[0]) {
		t.Errorf("Expected \"%s\" but get \"%s\".\n",
			data.expectedProperty, available[0])
	}

	// Free results hash
	err = results.Free()
	if err != nil {
		t.Error("Failed to free results hash object.")
	}

	// Free resource manager
	freeResourceManager(manager, t)
}

// This tests property name is returned correctly from a
// returned property index.
func TestPropertyName(t *testing.T) {
	properties := make(map[string]bool)
	properties["ismobile"] = true
	properties["browsername"] = true
	propertiesString := ""
	for k, _ := range properties {
		if propertiesString != "" {
			propertiesString += ","
		}
		propertiesString += k
	}

	manager := initResourceManager(propertiesString, t)
	// Create results hash
	results, err := NewResultsHash(manager, 1, 0)
	if err != nil {
		t.Error("Failed to create results hash object.")
	}

	available, err := results.AvailableProperties()
	if err != nil {
		t.Error("Failed to obtain available properties.")
	}

	if len(available) != len(properties) {
		t.Errorf("Expected to have %d properties but get %d.\n",
			len(properties),
			len(available))
	}

	for i, a := range available {
		name, err := results.PropertyName(i)
		if err != nil {
			t.Errorf("Failed to obtain property name for index %d.\n", i)
		}

		if !strings.EqualFold(a, name) {
			t.Errorf("Expected property \"%s\" but get \"%s\".\n",
				a, name)
		}

		if v := properties[strings.ToLower(name)]; v != true {
			t.Errorf("Property \"%s\" was not required or has already been "+
				"checked.\n", name)
		} else {
			properties[name] = false
		}
	}

	// Free results hash
	err = results.Free()
	if err != nil {
		t.Error("Failed to free results hash object.")
	}

	// Free resource manager
	freeResourceManager(manager, t)
}

// This tests a required property index is returned correctly
// from a given name
func TestRequiredPropertyIndexFromName(t *testing.T) {
	properties := make(map[string]bool)
	properties["ismobile"] = true
	properties["browsername"] = true
	propertiesString := ""
	for k, _ := range properties {
		if propertiesString != "" {
			propertiesString += ","
		}
		propertiesString += k
	}

	manager := initResourceManager(propertiesString, t)
	// Create results hash
	results, err := NewResultsHash(manager, 1, 0)
	if err != nil {
		t.Error("Failed to create results hash object.")
	}

	available, err := results.AvailableProperties()
	if err != nil {
		t.Error("Failed to obtain available properties.")
	}

	if len(available) != len(properties) {
		t.Errorf("Expected to have %d properties but get %d.\n",
			len(properties),
			len(available))
	}

	for i, name := range available {
		index, err := results.RequiredPropertyIndexFromName(name)
		if err != nil {
			t.Errorf("Failed to obtain property name for index %d.\n", i)
		}

		if index != i {
			t.Errorf("Expected index %d for property \"%s\" but get %d.\n ",
				i, name, index)
		} else {
			properties[name] = false
		}
	}

	// Free results hash
	err = results.Free()
	if err != nil {
		t.Error("Failed to free results hash object.")
	}

	// Free resource manager
	freeResourceManager(manager, t)
}
