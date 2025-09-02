package connections

import (
	"encoding/json"
	"log"
	"os"
)

type BitRateList struct {
	BitRates []BitRate `json:"bitrates"`
}

func (b *BitRateList) GetBitRates() []BitRate {
	return b.BitRates
}

func (b *BitRateList) SetBitRates(bitRates []BitRate) {
	b.BitRates = bitRates
}

type BitRate struct {
	Modulation string  `json:"modulation"`
	Slots      []Slots `json:"slots"`
	Reachs     []Reach `json:"reach"`
}

func (b *BitRate) GetModulation() string {
	return b.Modulation
}

func (b *BitRate) SetModulation(modulation string) {
	b.Modulation = modulation
}

func (b *BitRate) GetSlots() []Slots {
	return b.Slots
}

func (b *BitRate) SetSlots(slots []Slots) {
	b.Slots = slots
}

func (b *BitRate) GetReachs() []Reach {
	return b.Reachs
}

func (b *BitRate) SetReachs(reachs []Reach) {
	b.Reachs = reachs
}

type Slots struct {
	Gigabits string `json:"gigabits"`
	Slots    int    `json:"slots"`
}

func (s *Slots) GetGigabits() string {
	return s.Gigabits
}

func (s *Slots) SetGigabits(gigabits string) {
	s.Gigabits = gigabits
}

func (s *Slots) GetSlots() int {
	return s.Slots
}

func (s *Slots) SetSlots(slots int) {
	s.Slots = slots
}

type Reach struct {
	NumberOfBands int            `json:"number_of_bands"`
	Reach         []ReachPerBand `json:"reach"`
}

func (r *Reach) GetNumberOfBands() int {
	return r.NumberOfBands
}

func (r *Reach) SetNumberOfBands(numberOfBands int) {
	r.NumberOfBands = numberOfBands
}

func (r *Reach) GetReach() []ReachPerBand {
	return r.Reach
}

func (r *Reach) SetReach(reach []ReachPerBand) {
	r.Reach = reach
}

type ReachPerBand struct {
	Band  string `json:"band"`
	Reach int    `json:"reach"`
}

func (r *ReachPerBand) GetBand() string {
	return r.Band
}

func (r *ReachPerBand) SetBand(band string) {
	r.Band = band
}

func (r *ReachPerBand) GetReach() int {
	return r.Reach
}

func (r *ReachPerBand) SetReach(reach int) {
	r.Reach = reach
}

func ReadBitRateFile(bitRatePath string) (BitRateList, error) {
	dataBytesBitrate, err := os.ReadFile(bitRatePath)

	if err != nil {
		log.Fatalf("Error opening file: %v", err)
		return BitRateList{}, err
	}

	var bitrate BitRateList

	err = json.Unmarshal(dataBytesBitrate, &bitrate)

	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
		return BitRateList{}, err
	}

	return bitrate, nil
}

func TrasnformIntToModulation(modulation int) string {
	switch modulation {
	case 0:
		return "BPSK"
	case 1:
		return "QPSK"
	case 2:
		return "8-QAM"
	case 3:
		return "16-QAM"
	default:
		log.Fatalf("Invalid modulation type: %d", modulation)
		return "BPSK" // Default case, should not be reached
	}
}
