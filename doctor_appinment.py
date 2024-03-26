class Doctor:
    def __init__(self, name, speciality, rating=4):
        self.name = name
        self.speciality = speciality
        self.availability = []
        self.rating = rating

    def declare_availability(self, slots):
        self.availability.extend(slots)


class Patient:
    def __init__(self, name):
        self.name = name
        self.booked_appointments = {}

    def book_appointment(self, booking_id, slot, doctor_name):
        self.booked_appointments[booking_id] = {"slot": slot, "doctor": doctor_name}

    def cancel_appointment(self, booking_id):
        if booking_id in self.booked_appointments:
            del self.booked_appointments[booking_id]
            return True
        return False


class AppointmentSystem:
    def __init__(self):
        self.doctors = []
        self.patients = []
        self.bookings = {}
        self.booking_id = 1000

    def register_doctor(self, name, speciality, rating=4):
        doctor = Doctor(name, speciality, rating)
        self.doctors.append(doctor)
        print(f"Welcome Dr. {name} !!")

    def mark_doc_avail(self, doctor_name, slots):
        doctor = self.find_doctor_by_name(doctor_name)
        if not doctor:
            print("Doctor not found!")
            return

        for slot in slots:
            if not self.is_valid_slot(slot):
                print(f"Invalid slot format for Dr. {doctor_name}. Slots must be in format 'hh:mm-hh:mm'")
                return

            start, end = slot.split("-")
            start_time = int(start[:2]) * 60 + int(start[3:])
            end_time = int(end[:2]) * 60 + int(end[3:])
            if end_time - start_time != 60:
                print(f"Invalid slot duration for Dr. {doctor_name}. Slots must be exactly 60 minutes long.")
                return

            doctor.declare_availability(slots)
        print("Done Doc!")

    def rank_doctors_by_rating(self, speciality):
        ranked_doctors = {}
        for doctor in self.doctors:
            if doctor.speciality == speciality:
                ranked_doctors[doctor.name] = doctor.rating
        return dict(sorted(ranked_doctors.items(), key=lambda item: item[1], reverse=True))

    def show_avail_byspeciality(self, speciality, ranking_strategy='start_time'):
        print(speciality)
        if ranking_strategy == 'rating':
            ranked_doctors = self.rank_doctors_by_rating(speciality)
            for doctor_name, _ in ranked_doctors.items():
                self.show_avail_for_doctor(doctor_name)
        else:
            for doctor in self.doctors:
                if doctor.speciality == speciality:
                    self.show_avail_for_doctor(doctor.name)

    def show_avail_for_doctor(self, doctor_name):
        doctor = self.find_doctor_by_name(doctor_name)
        for slot in doctor.availability:
            print(f"Dr. {doctor.name}: ({slot})")
        print()

    def register_patient(self, name):
        patient = Patient(name)
        self.patients.append(patient)
        print(f"{name} registered successfully.")

    def book_appointment(self, patient_name, doctor_name, slot, waitlist=False):
        patient = self.find_patient_by_name(patient_name)
        doctor = self.find_doctor_by_name(doctor_name)

        if not patient or not doctor:
            print("Patient or Doctor not found!")
            return

        if slot not in doctor.availability:
            print("Slot not available for booking!")
            return

        if not waitlist and self.is_slot_booked(slot):
            print("Slot already booked!")
            return

        booking_id = self.booking_id + 1
        self.booking_id += 1
        self.bookings[booking_id] = {"patient": patient_name, "doctor": doctor_name, "slot": slot}
        patient.book_appointment(booking_id, slot, doctor_name)
        doctor.availability.remove(slot)
        print(f"Booked. Booking id: {booking_id}")

    def cancel_booking_id(self, booking_id):
        if booking_id not in self.bookings:
            print("Booking not found!")
            return

        slot = self.bookings[booking_id]["slot"]
        doctor_name = self.bookings[booking_id]["doctor"]
        del self.bookings[booking_id]

        for patient in self.patients:
            if patient.cancel_appointment(booking_id):
                break

        self.find_doctor_by_name(doctor_name).availability.append(slot)
        print("Booking Cancelled")

    def show_appointments_booked(self, patient_name):
        patient = self.find_patient_by_name(patient_name)
        if not patient:
            print("Patient not found!")
            return

        for booking_id, slot_info in patient.booked_appointments.items():
            print(f"Booking id: {booking_id}, Dr {slot_info['doctor']} {slot_info['slot']}")

    def find_doctor_by_name(self, name):
        for doctor in self.doctors:
            if doctor.name == name:
                return doctor
        return None

    def find_patient_by_name(self, name):
        for patient in self.patients:
            if patient.name == name:
                return patient
        return None

    def is_valid_slot(self, slot):
        return len(slot) == 11 and slot[2] == ":" and slot[5] == "-" and slot[8] == ":"

    def is_slot_booked(self, slot):
        for booking in self.bookings.values():
            if booking["slot"] == slot:
                return True
        return False


system = AppointmentSystem()
print("Enter command (registerDoc, markDocAvail, showAvailByspeciality, registerPatient, bookAppointment, cancelBookingId, showAppointmentsBooked, exit):")
while True:
    command = input()
    if command == 'exit':
        break

    args = command.split(" ")
    action = args[0]

    if action == 'registerDoc':
        name, speciality = args[1], args[2]
        system.register_doctor(name, speciality)
    elif action == 'markDocAvail':
        name, slots = args[1], args[2:]
        system.mark_doc_avail(name, slots)
    elif action == 'showAvailByspeciality':
        speciality = args[1]
        system.show_avail_byspeciality(speciality)
    elif action == 'registerPatient':
        name = args[1]
        system.register_patient(name)
    elif action == 'bookAppointment':
        patient_name, doctor_name, slot = args[1], args[2], args[3]
        waitlist = True if len(args) > 4 and args[4] == 'true' else False
        system.book_appointment(patient_name, doctor_name, slot, waitlist)
    elif action == 'cancelBookingId':
        booking_id = args[1]
        system.cancel_booking_id(booking_id)
    elif action == 'showAppointmentsBooked':
        patient_name = args[1]
        system.show_appointments_booked(patient_name)
    else:
        print("Invalid command!")
