package data

type DetailedArtist struct {
	Artist
	LocationsData  []string
	DatesData      []string
	DatesLocations map[string][]string
}

func WeaveData(
	artists []Artist,
	locationsIndex []Location,
	datesIndex []Date,
	relationsIndex []Relation,
) map[int]DetailedArtist {

	// 1. Create the master map
	masterMap := make(map[int]DetailedArtist, len(artists))

	// 2. Add artists
	for _, artist := range artists {
		masterMap[artist.ID] = DetailedArtist{
			Artist:         artist,
			LocationsData:  []string{},
			DatesData:      []string{},
			DatesLocations: make(map[string][]string),
		}
	}

	// 3. Add locations
	for _, l := range locationsIndex {
		detailed, exists := masterMap[l.ID]
		if !exists {
			continue
		}

		detailed.LocationsData = l.Locations
		masterMap[l.ID] = detailed
	}

	// 4. Add dates
	for _, d := range datesIndex {
		detailed, exists := masterMap[d.ID]
		if !exists {
			continue
		}

		detailed.DatesData = d.Dates
		masterMap[d.ID] = detailed
	}

	// 5. Add relations
	for _, r := range relationsIndex {
		detailed, exists := masterMap[r.ID]
		if !exists {
			continue
		}

		detailed.DatesLocations = r.DatesLocations
		masterMap[r.ID] = detailed
	}

	return masterMap
}
