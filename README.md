Project 5
================================================

### What

Healthcare data

### Why/How

Today, doctors, nurses & health professionals are limited in the level of care they can provide, because they can not view complete, accurate health records. Doctors will be able to publish medical records safely on the blockchain and only authorized people will be able to access the data.
Medical institutions will be able to securely share data, no medical records will be lost. Test results will come on time and the probability of correct diagnosis and effective treatment will be higher. Blockchain is immutable, so nobody will be able to fake or lose data. Blockchain is decentralized, so the information will be stored not only in one datacenter, in this case - in different medical institutions.

### Functionalities/Success

Doctor can add information about patient
Information will be available to view
Information will be non changeable
Patient will be able to see itâ€™s data by entering personal code

### Doctors Application
Doc's data: Public Key, Private Key, DocID
Doctor knows Public Keys, Private Key and PatID of patients
#### Functionalities
##### Add data
* Type patient ID and Patient Info into webPage form
* If Patient ID doesn't exist, the new patient is added to patientList (new Public/Private keys are created)
* Data from webPage will be added in format <PatID, PatInfo> for user's DocID

##### Send data to miners
Data format: <DocID, <PatID, PatInfo> >
* Encrypt PatientInfo with patient's Private Key
* Hash <PatID, [PatInfo]PK>
* Sign H<PatID, [PatInfo]PKpat> with doc's Private key
* Send to miner <DocID, [H<PatID, [PatInfo]PKpat>]PKdoc>

##### Read Data
Doctor reads only his/her data (via ID)
Doctors can read data only from canonical chain
* Search all values with key == DoctorID in canonical blocks
* Decrypt [H<PatID, [PatInfo]PKpat>]PKdoc with Doc's Private Key
* Dehash [H<PatID, [PatInfo]PKpat>]
* Compare Hash with new generated Hash of <PatID, [PatInfo]PKpat> (Verify that data hasn't been changed)
* Decrypt [PatInfo]PKpat by Pat's Public Keys for each PatientID
* Add to result
* Finished result is printed to web page in form:
```
Patients results
Block № 2
Timestamp = 1556758412624725900
Patient ID: 12345, Patient Data = In order to avoid heart disease, it is important to eat a healthy diet and take regular exercise.
Patient ID: 23456, Patient Data = Each year millions of people take time off work due to influenza in the winter months.
Patient ID: 34567, Patient Data = Viral diseases are often harder to treat than bacterial infections.

Block № 1
Timestamp = 1556758399042040600
Patient ID: 12345, Patient Data = You have to be very careful how much sugar you eat when you have diabetes.
```


### Patients Application
Patient's data: Public Key, Private Key, PatID
Patient knows Public Keys and DocID of doctors
#### Functionalities
##### Read his/her data
Patient reads only his/her data (via ID)
Patient can read data only from canonical chain
* Enter DocID into WebPage
* Search all values with key == DoctorID in canonical blocks
* Decrypt [H<PatID, [PatInfo]PKpat>]PKdoc with Doc's Private Key
* Dehash [H<PatID, [PatInfo]PKpat>]
* Compare Hash with new generated Hash of <PatID, [PatInfo]PKpat> (Verify that data hasn't been changed)
* Search all key-values with key == PatID
* Decrypt [PatInfo]PKpat by Pat's Private Key
* Add to result
* Finished result is printed to web page in form:

```
Patient ID 12345 results

Timestamp = 1556758412624725900
Patient ID: 12345, Patient Data = In order to avoid heart disease, it is important to eat a healthy diet and take regular exercise.

Timestamp = 1556758399042040600
Patient ID: 12345, Patient Data = You have to be very careful how much sugar you eat when you have diabetes.
```


### API Functionalities for DoctorsApp

GET /patients
Web Page with information about all patients in blockchain. Data shows only canonical chain without forks.
The data contains block number, timestamp, patient ID and patient information.

GET /add
WebPage to write PatientID - PatientInfo
POST /add
Add information by doctor, send to miners

### API Functionalities for PatientsApp

GET /patient
WebPage to write DocID
POST /patient
See information about a particular patient

### Use Case Diagram

![alt text](https://github.com/BettyRain/blockchain_project5/blob/master/Project5UseCaseDiagram.jpeg)

### Timeline and milestones

| N%   | Milestones         | Completion Date |
| :-------: |:-------------: | :-------------|
| 1 | Add patient data by Doc's app | April 28nd |
| 2 | View data from blocks by Doc's app | April 28th |
| 3 | Make interface to adding data (web page) | April 30th |
| 4 | Checkpoint | May 1st |
| 5 | Data confidentiality | May 8th |
| 6 | Data integrity | May 10th |
| 7 | View data by patients (Patient's App) | May 10th |

### References
Idea was found at [blockchain ideas](https://www.connectbit.com/blockchain-applications/)