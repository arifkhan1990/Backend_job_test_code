class Doctor {
  constructor(name, speciality, rating) {
    this.name = name;
    this.speciality = speciality;
    this.availability = [];
    this.rating = rating;
  }

  declareAvailability(slots) {
    slots.forEach((slot) => {
      this.availability.push(slot);
    });
  }
}

class Patient {
  constructor(name) {
    this.name = name;
    this.bookedAppointments = {};
  }

  bookAppointment(bookingId, slot, doctorName) {
    this.bookedAppointments[bookingId] = { slot, doctor: doctorName };
  }

  cancelAppointment(bookingId) {
    if (this.bookedAppointments.hasOwnProperty(bookingId)) {
      delete this.bookedAppointments[bookingId];
      return true;
    }
    return false;
  }
}

class AppointmentSystem {
  constructor() {
    this.doctors = [];
    this.patients = [];
    this.bookings = {};
    this.bookingId = 1000;
  }

  registerDoctor(name, speciality, rating = 4) {
    const doctor = new Doctor(name, speciality, rating);
    this.doctors.push(doctor);
    console.log(`Welcome Dr. ${name} !!`);
  }

  markDocAvail(doctorName, slots) {
    const doctor = this.findDoctorByName(doctorName);
    if (!doctor) {
      console.log("Doctor not found!");
      return;
    }
    slots.forEach((slot) => {
      slot = slot.replace(",", "");

      if (!this.isValidSlot(slot)) {
        console.log(
          `Invalid slot format for Dr. ${doctorName}. Slots must be in format 'hh:mm-hh:mm'`
        );
        return;
      }

      const [start, end] = slot.split("-");
      const startTime = new Date(`2022-01-01T${start}`);
      const endTime = new Date(`2022-01-01T${end}`);
      if (endTime - startTime !== 3600000) {
        console.log(
          `Invalid slot duration for Dr. ${doctorName}. Slots must be exactly 60 minutes long.`
        );
        return;
      }
      this.addSlotToDoctor(doctor, slot);
    });

    console.log("Done Doc!");
  }

  isValidSlot(slot) {
    const pattern =
      /^(\d{1,2}|\d{2}):(\d{1,2}|\d{2})-(\d{1,2}|\d{2}):(\d{1,2}|\d{2})$/;
    return pattern.test(slot);
  }

  addSlotToDoctor(doctor, slot) {
    doctor.availability.push(slot);
  }

  rankDoctorsByRating(speciality) {
    const rankedDoctors = {};
    this.doctors.forEach((doctor) => {
      if (doctor.speciality === speciality) {
        rankedDoctors[doctor.name] = doctor.rating;
      }
    });
    return Object.fromEntries(
      Object.entries(rankedDoctors).sort((a, b) => b[1] - a[1])
    );
  }

  showAvailByspeciality(speciality, rankingStrategy = "start_time") {
    console.log(speciality);
    if (rankingStrategy === "rating") {
      const rankedDoctors = this.rankDoctorsByRating(speciality);
      Object.keys(rankedDoctors).forEach((doctorName) => {
        this.showAvailForDoctor(doctorName, speciality);
      });
    } else {
      this.doctors.forEach((doctor) => {
        if (doctor.speciality === speciality) {
          this.showAvailForDoctor(doctor.name, speciality);
        }
      });
    }
  }

  showAvailForDoctor(doctorName, speciality) {
    const doctor = this.findDoctorByName(doctorName);
    doctor.availability.forEach((slot) => {
      console.log(`Dr. ${doctor.name}: (${slot})`);
    });
    console.log("");
  }

  registerPatient(name) {
    const patient = new Patient(name);
    this.patients.push(patient);
    console.log(`${name} registered successfully.`);
  }

  bookAppointment(patientName, doctorName, slot, waitlist = false) {
    const patient = this.findPatientByName(patientName);
    const doctor = this.findDoctorByName(doctorName);

    if (!patient || !doctor) {
      console.log("Patient or Doctor not found!");
      return;
    }

    if (!doctor.availability.includes(slot)) {
      console.log("Slot not available for booking!");
      return;
    }

    if (!waitlist && this.isSlotBooked(slot)) {
      console.log("Slot already booked!");
      return;
    }

    const bookingId = ++this.bookingId;
    this.bookings[bookingId] = {
      patient: patientName,
      doctor: doctorName,
      slot,
    };
    patient.bookAppointment(bookingId, slot, doctorName);
    this.removeSlotFromDoctor(doctor, slot);
    console.log(`Booked. Booking id: ${bookingId}`);
  }

  removeSlotFromDoctor(doctor, slotToRemove) {
    doctor.availability = doctor.availability.filter(
      (slot) => slot !== slotToRemove
    );
  }

  cancelBookingId(bookingId) {
    if (!this.bookings.hasOwnProperty(bookingId)) {
      console.log("Booking not found!");
      return;
    }

    const slot = this.bookings[bookingId]["slot"];
    const doctorName = this.bookings[bookingId]["doctor"];
    delete this.bookings[bookingId];
    this.patients.forEach((patient) => {
      if (patient.cancelAppointment(bookingId)) {
        return; // Cancel only once for the given booking ID
      }
    });

    // Add cancelled slot back to doctor's availability
    const doctor = this.findDoctorByName(doctorName);
    this.addSlotToDoctor(doctor, slot);
    console.log("Booking Cancelled");
  }

  showAppointmentsBooked(patientName) {
    const patient = this.findPatientByName(patientName);
    if (!patient) {
      console.log("Patient not found!");
      return;
    }

    Object.values(patient.bookedAppointments).forEach((slot) => {
      const bookingId = Object.keys(this.bookings).find(
        (key) => this.bookings[key]["slot"] === slot.slot
      );
      const doctor = this.bookings[bookingId]["doctor"];
      console.log(`Booking id: ${bookingId}, Dr ${doctor} ${slot.slot}`);
    });
  }

  findDoctorByName(name) {
    return this.doctors.find((doctor) => doctor.name === name);
  }

  findPatientByName(name) {
    return this.patients.find((patient) => patient.name === name);
  }

  isSlotBooked(slot) {
    return Object.values(this.bookings).some(
      (booking) => booking.slot === slot
    );
  }
}

const system = new AppointmentSystem();
console.log(
  "Enter command (registerDoc, markDocAvail, showAvailByspeciality, registerPatient, bookAppointment, cancelBookingId, showAppointmentsBooked, exit): "
);
const readline = require("readline");
const rl = readline.createInterface({
  input: process.stdin,
  output: process.stdout,
});

rl.on("line", (input) => {
  const command = input.trim();
  if (command === "exit") {
    rl.close();
    return;
  }

  const args = command.split(" ");
  const action = args[0];

  switch (action) {
    case "registerDoc":
      const docName = args[1];
      const speciality = args[2];
      system.registerDoctor(docName, speciality);
      break;
    case "markDocAvail":
      const doctorName = args[1];
      const slots = args.slice(2);
      system.markDocAvail(doctorName, slots);
      break;
    case "showAvailByspeciality":
      const specialityToShow = args[1];
      system.showAvailByspeciality(specialityToShow);
      break;
    case "registerPatient":
      const patientName = args[1];
      system.registerPatient(patientName);
      break;
    case "bookAppointment":
      const patient = args[1];
      const docNameForBooking = args[2];
      const slotForBooking = args[3];
      const isWaitlist = args[4] === "true";
      system.bookAppointment(
        patient,
        docNameForBooking,
        slotForBooking,
        isWaitlist
      );
      break;
    case "cancelBookingId":
      const bookingID = args[1];
      system.cancelBookingId(bookingID);
      break;
    case "showAppointmentsBooked":
      const patientNameForAppointments = args[1];
      system.showAppointmentsBooked(patientNameForAppointments);
      break;
    default:
      console.log("Invalid command!");
      break;
  }
});

rl.on("close", () => {
  console.log("Exiting...");
});
