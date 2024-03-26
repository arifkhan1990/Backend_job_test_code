<?php

class Doctor
{
    public $name;
    public $speciality;
    public $availability = [];
    public $rating;

    public function __construct($name, $speciality, $rating)
    {
        $this->name = $name;
        $this->speciality = $speciality;
        $this->rating = $rating;
    }

    public function declareAvailability($slots)
    {
        foreach ($slots as $slot) {
            $this->availability[] = $slot;
        }
    }
}

class Patient
{
    public $name;
    public $bookedAppointments = [];

    public function __construct($name)
    {
        $this->name = $name;
    }

    public function bookAppointment($bookingId, $slot, $doctorName)
    {
        $this->bookedAppointments[$bookingId] = ["slot" => $slot, "doctor" => $doctorName];
    }

    public function cancelAppointment($bookingId)
    {
        if (isset($this->bookedAppointments[$bookingId])) {
            unset($this->bookedAppointments[$bookingId]);
            return true;
        }
        return false;
    }
}

class AppointmentSystem
{
    public $doctors = [];
    public $patients = [];
    public $bookings = [];
    public $bookingId = 1000;

    public function registerDoctor($name, $speciality, $rating = 4)
    {
        $doctor = new Doctor($name, $speciality, $rating);
        $this->doctors[] = $doctor;
        echo "Welcome Dr. $name !!\n";
    }

    public function markDocAvail($doctorName, $slots)
    {
        $doctor = $this->findDoctorByName($doctorName);
        if (!$doctor) {
            echo "Doctor not found!\n";
            return;
        }
        foreach ($slots as $slot) {
            $slot = rtrim($slot, ',');

            if (!$this->isValidSlot($slot)) {
                echo "Invalid slot format for Dr. $doctorName. Slots must be in format 'hh:mm-hh:mm'\n";
                return;
            }

            list($start, $end) = explode('-', $slot);
            $startTime = strtotime($start);
            $endTime = strtotime($end);
            if (($endTime - $startTime) !== 3600) {
                echo "Invalid slot duration for Dr. $doctorName. Slots must be exactly 60 minutes long.\n";
                return;
            }
            $this->addSlotToDoctor($doctor, $slot);
        }

        echo "Done Doc!\n";
    }

    private function isValidSlot($slot)
    {
        $pattern = '/^(\d{1,2}|\d{2}):(\d{1,2}|\d{2})-(\d{1,2}|\d{2}):(\d{1,2}|\d{2})$/';
        return preg_match($pattern, $slot);
    }

    private function addSlotToDoctor($doctor, $slot)
    {
        list($startTime, $endTime) = explode('-', $slot);
        $doctor->availability[] = $slot;
    }

    public function rankDoctorsByRating($speciality)
    {
        $rankedDoctors = [];
        foreach ($this->doctors as $doctor) {
            if ($doctor->speciality === $speciality) {
                $rankedDoctors[$doctor->name] = $doctor->rating;
            }
        }
        arsort($rankedDoctors);
        return $rankedDoctors;
    }
    // Display slots in a ranked fashion
    public function showAvailByspeciality($speciality, $rankingStrategy = 'start_time')
    {
        echo $speciality;
        if ($rankingStrategy === 'rating') {
            $rankedDoctors = $this->rankDoctorsByRating($speciality);
            foreach ($rankedDoctors as $doctorName => $rating) {
                $this->showAvailForDoctor($doctorName, $speciality);
            }
        } else {
            foreach ($this->doctors as $doctor) {
                if ($doctor->speciality === $speciality) {
                    $this->showAvailForDoctor($doctor->name, $speciality);
                }
            }
        }
    }

    private function showAvailForDoctor($doctorName, $speciality)
    {
        $doctor = $this->findDoctorByName($doctorName);
        foreach ($doctor->availability as $slot) {
            echo "Dr. $doctor->name: ";
            echo "($slot) \n";
        }
        echo "\n";
    }

    public function registerPatient($name)
    {
        $patient = new Patient($name);
        $this->patients[] = $patient;
        echo "$name registered successfully.\n";
    }

    public function bookAppointment($patientName, $doctorName, $slot, $waitlist = false)
    {
        $patient = $this->findPatientByName($patientName);
        $doctor = $this->findDoctorByName($doctorName);

        if (!$patient || !$doctor) {
            echo "Patient or Doctor not found!\n";
            return;
        }

        if (!in_array($slot, $doctor->availability)) {
            echo "Slot not available for booking!\n";
            return;
        }

        if (!$waitlist && $this->isSlotBooked($slot)) {
            echo "Slot already booked!\n";
            return;
        }

        $bookingId = ++$this->bookingId;
        $this->bookings[$bookingId] = ["patient" => $patientName, "doctor" => $doctorName, "slot" => $slot];
        $patient->bookAppointment($bookingId, $slot, $doctorName);
        $this->removeSlotFromDoctor($doctor, $slot);
        echo "Booked. Booking id: $bookingId\n";
    }

    private function removeSlotFromDoctor($doctor, $slotToRemove)
    {
        $updatedAvailability = [];
        foreach ($doctor->availability as $slot) {
            if ($slot !== $slotToRemove) {
                $updatedAvailability[] = $slot;
            }
        }
        $doctor->availability = $updatedAvailability;
    }

    public function cancelBookingId($bookingId)
    {
        if (!isset($this->bookings[$bookingId])) {
            echo "Booking not found!\n";
            return;
        }

        $slot = $this->bookings[$bookingId]["slot"];
        $doctorName = $this->bookings[$bookingId]["doctor"];
        unset($this->bookings[$bookingId]);
        foreach ($this->patients as $patient) {
            if ($patient->cancelAppointment($bookingId)) {
                break; // Cancel only once for the given booking ID
            }
        }

        // Add cancelled slot back to doctor's availability
        $this->addSlotToDoctor($this->findDoctorByName($doctorName), $slot);
        echo "Booking Cancelled\n";
    }

    public function showAppointmentsBooked($patientName)
    {
        $patient = $this->findPatientByName($patientName);
        if (!$patient) {
            echo "Patient not found!\n";
            return;
        }

        foreach ($patient->bookedAppointments as $slot) {
            $bookingId = array_search($slot, array_column($this->bookings, 'slot'));
            $doctor = $this->bookings[$bookingId]["doctor"];
            echo "Booking id: $bookingId, Dr $doctor $slot\n";
        }
    }

    private function findDoctorByName($name)
    {
        foreach ($this->doctors as $doctor) {
            if ($doctor->name === $name) {

                return $doctor;
            }
        }
        return null;
    }

    private function findPatientByName($name)
    {
        foreach ($this->patients as $patient) {
            if ($patient->name === $name) {
                return $patient;
            }
        }
        return null;
    }

    private function isSlotBooked($slot)
    {
        foreach ($this->bookings as $booking) {
            if ($booking["slot"] === $slot) {
                return true;
            }
        }
        return false;
    }
}

$system = new AppointmentSystem();
echo "Enter command (registerDoc, markDocAvail, showAvailByspeciality, registerPatient, bookAppointment, cancelBookingId, showAppointmentsBooked, exit): \n";
while (true) {
    $command = readline();
    if ($command === 'exit') {
        break;
    }

    $args = explode(" ", $command);
    $action = $args[0];

    switch ($action) {
        case 'registerDoc':
            $name = $args[1];
            $speciality = $args[2];
            $system->registerDoctor($name, $speciality);
            break;
        case 'markDocAvail':
            $name = $args[1];
            $slots = array_slice($args, 2);
            $system->markDocAvail($name, $slots);
            break;
        case 'showAvailByspeciality':
            $speciality = $args[1];
            $system->showAvailByspeciality($speciality);
            break;
        case 'registerPatient':
            $name = $args[1];
            $system->registerPatient($name);
            break;
        case 'bookAppointment':
            $patientName = $args[1];
            $doctorName = $args[2];
            $slot = $args[3];
            $waitlist = isset($args[4]) && $args[4] === 'true';
            $system->bookAppointment($patientName, $doctorName, $slot, $waitlist);
            break;
        case 'cancelBookingId':
            $bookingId = $args[1];
            $system->cancelBookingId($bookingId);
            break;
        case 'showAppointmentsBooked':
            $patientName = $args[1];
            $system->showAppointmentsBooked($patientName);
            break;
        default:
            echo "Invalid command!\n";
            break;
    }
}
