package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Doctor struct {
	Name         string
	Speciality   string
	Availability []string
	Rating       int
}

type Patient struct {
	Name               string
	BookedAppointments map[string]map[string]string
}

type AppointmentSystem struct {
	Doctors   []*Doctor
	Patients  []*Patient
	Bookings  map[string]map[string]string
	BookingID int
}

func NewAppointmentSystem() *AppointmentSystem {
	return &AppointmentSystem{
		Doctors:   []*Doctor{},
		Patients:  []*Patient{},
		Bookings:  map[string]map[string]string{},
		BookingID: 1000,
	}
}

func NewDoctor(name, speciality string, rating int) *Doctor {
	return &Doctor{
		Name:         name,
		Speciality:   speciality,
		Availability: []string{},
		Rating:       rating,
	}
}

func NewPatient(name string) *Patient {
	return &Patient{
		Name:               name,
		BookedAppointments: map[string]map[string]string{},
	}
}

func (d *Doctor) DeclareAvailability(slots []string) {
	d.Availability = append(d.Availability, slots...)
}

func (p *Patient) BookAppointment(bookingID, slot, doctorName string) {
	if p.BookedAppointments == nil {
		p.BookedAppointments = map[string]map[string]string{}
	}
	p.BookedAppointments[bookingID] = map[string]string{"slot": slot, "doctor": doctorName}
}

func (p *Patient) CancelAppointment(bookingID string) bool {
	if _, exists := p.BookedAppointments[bookingID]; exists {
		delete(p.BookedAppointments, bookingID)
		return true
	}
	return false
}

func (as *AppointmentSystem) RegisterDoctor(name, speciality string, rating int) {
	doctor := NewDoctor(name, speciality, rating)
	as.Doctors = append(as.Doctors, doctor)
	fmt.Printf("Welcome Dr. %s !!\n", name)
}

func (as *AppointmentSystem) MarkDocAvail(doctorName string, slots []string) {
	doctor := as.FindDoctorByName(doctorName)
	if doctor == nil {
		fmt.Println("Doctor not found!")
		return
	}

	for _, slot := range slots {
		if !as.IsValidSlot(slot) {
			fmt.Printf("Invalid slot format for Dr. %s. Slots must be in format 'hh:mm-hh:mm'\n", doctorName)
			return
		}

		startEnd := strings.Split(slot, "-")
		if len(startEnd) != 2 {
			fmt.Printf("Invalid slot format for Dr. %s. Slots must be in format 'hh:mm-hh:mm'\n", doctorName)
			return
		}
		doctor.DeclareAvailability(slots)
	}
	fmt.Println("Done Doc!")
}

func (as *AppointmentSystem) RankDoctorsByRating(speciality string) map[string]int {
	rankedDoctors := map[string]int{}
	for _, doctor := range as.Doctors {
		if doctor.Speciality == speciality {
			rankedDoctors[doctor.Name] = doctor.Rating
		}
	}
	return rankedDoctors
}

func (as *AppointmentSystem) ShowAvailBySpeciality(speciality string, rankingStrategy string) {
	fmt.Println(speciality)
	if rankingStrategy == "rating" {
		rankedDoctors := as.RankDoctorsByRating(speciality)
		for doctorName, _ := range rankedDoctors {
			as.ShowAvailForDoctor(doctorName)
		}
	} else {
		for _, doctor := range as.Doctors {
			if doctor.Speciality == speciality {
				as.ShowAvailForDoctor(doctor.Name)
			}
		}
	}
}

func (as *AppointmentSystem) ShowAvailForDoctor(doctorName string) {
	doctor := as.FindDoctorByName(doctorName)
	if doctor != nil {
		for _, slot := range doctor.Availability {
			fmt.Printf("Dr. %s: (%s)\n", doctor.Name, slot)
		}
		fmt.Println()
	}
}

func (as *AppointmentSystem) RegisterPatient(name string) {
	patient := NewPatient(name)
	as.Patients = append(as.Patients, patient)
	fmt.Printf("%s registered successfully.\n", name)
}

func (as *AppointmentSystem) BookAppointment(patientName, doctorName, slot string, waitlist bool) {
	patient := as.FindPatientByName(patientName)
	doctor := as.FindDoctorByName(doctorName)

	if patient == nil || doctor == nil {
		fmt.Println("Patient or Doctor not found!")
		return
	}

	if !waitlist && as.IsSlotBooked(slot) {
		fmt.Println("Slot already booked!")
		return
	}

	if !as.IsSlotAvailable(doctor, slot) {
		fmt.Println("Slot not available for booking!")
		return
	}

	bookingID := strconv.Itoa(as.BookingID + 1)
	as.BookingID++
	as.Bookings[bookingID] = map[string]string{"patient": patientName, "doctor": doctorName, "slot": slot}
	patient.BookAppointment(bookingID, slot, doctorName)
	doctor.Availability = RemoveString(doctor.Availability, slot)
	fmt.Printf("Booked. Booking id: %s\n", bookingID)
}

func (as *AppointmentSystem) CancelBookingID(bookingID string) {
	if _, exists := as.Bookings[bookingID]; !exists {
		fmt.Println("Booking not found!")
		return
	}

	slot := as.Bookings[bookingID]["slot"]
	doctorName := as.Bookings[bookingID]["doctor"]
	delete(as.Bookings, bookingID)

	for _, patient := range as.Patients {
		if patient.CancelAppointment(bookingID) {
			break
		}
	}

	as.FindDoctorByName(doctorName).Availability = append(as.FindDoctorByName(doctorName).Availability, slot)
	fmt.Println("Booking Cancelled")
}

func (as *AppointmentSystem) ShowAppointmentsBooked(patientName string) {
	patient := as.FindPatientByName(patientName)
	if patient == nil {
		fmt.Println("Patient not found!")
		return
	}

	for bookingID, slotInfo := range patient.BookedAppointments {
		fmt.Printf("Booking id: %s, Dr %s %s\n", bookingID, slotInfo["doctor"], slotInfo["slot"])
	}
}

func (as *AppointmentSystem) FindDoctorByName(name string) *Doctor {
	for _, doctor := range as.Doctors {
		if doctor.Name == name {
			return doctor
		}
	}
	return nil
}

func (as *AppointmentSystem) FindPatientByName(name string) *Patient {
	for _, patient := range as.Patients {
		if patient.Name == name {
			return patient
		}
	}
	return nil
}

func (as *AppointmentSystem) IsValidSlot(slot string) bool {
	return len(slot) == 11 && slot[2] == ':' && slot[5] == '-' && slot[8] == ':'
}

func (as *AppointmentSystem) IsSlotBooked(slot string) bool {
	for _, booking := range as.Bookings {
		if booking["slot"] == slot {
			return true
		}
	}
	return false
}

func (as *AppointmentSystem) IsSlotAvailable(doctor *Doctor, slot string) bool {
	for _, availableSlot := range doctor.Availability {
		if availableSlot == slot {
			return true
		}
	}
	return false
}

func RemoveString(slice []string, s string) []string {
	var result []string
	for _, str := range slice {
		if str != s {
			result = append(result, str)
		}
	}
	return result
}

func main() {
	system := NewAppointmentSystem()
	fmt.Println("Enter command (registerDoc, markDocAvail, showAvailByspeciality, registerPatient, bookAppointment, cancelBookingId, showAppointmentsBooked, exit):")
	var command string
	for {
		fmt.Scanln(&command)
		if command == "exit" {
			break
		}

		args := strings.Split(command, " ")
		action := args[0]

		switch action {
		case "registerDoc":
			name, speciality := args[1], args[2]
			rating := 4
			system.RegisterDoctor(name, speciality, rating)
		case "markDocAvail":
			name := args[1]
			slots := args[2:]
			system.MarkDocAvail(name, slots)
		case "showAvailByspeciality":
			speciality := args[1]
			rankingStrategy := "start_time"
			if len(args) > 2 {
				rankingStrategy = args[2]
			}
			system.ShowAvailBySpeciality(speciality, rankingStrategy)
		case "registerPatient":
			name := args[1]
			system.RegisterPatient(name)
		case "bookAppointment":
			patientName, doctorName, slot := args[1], args[2], args[3]
			waitlist := false
			if len(args) > 4 && args[4] == "true" {
				waitlist = true
			}
			system.BookAppointment(patientName, doctorName, slot, waitlist)
		case "cancelBookingId":
			bookingID := args[1]
			system.CancelBookingID(bookingID)
		case "showAppointmentsBooked":
			patientName := args[1]
			system.ShowAppointmentsBooked(patientName)
		default:
			fmt.Println("Invalid command!")
		}
		fmt.Println("Enter command:")
	}
}
