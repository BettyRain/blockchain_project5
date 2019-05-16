Project 5
================================================

### What

Healthcare data

### Why/How

Today, doctors, nurses & health professionals are limited in the level of care they can provide, because they can not view complete, accurate health records. Doctors will be able to publish medical records safely on the blockchain and only authorized people will be able to access the data. <br />
Medical institutions will be able to securely share data, no medical records will be lost. Test results will come on time and the probability of correct diagnosis and effective treatment will be higher. Blockchain is immutable, so nobody will be able to fake or lose data. Blockchain is decentralized, so the information will be stored not only in one datacenter, in this case - in different medical institutions.

### Functionalities/Success

* Doctor can add information about patient
* Information will be available to view
* Information will be non changeable
* Patient will be able to see his/her data

### Doctors Application
Doc's data: Public Key, Private Key, DocID <br />
Doctor knows Public Keys, Private Key and PatID of patients
#### Functionalities
##### Add data
* Type patient ID and Patient Info into webPage form
* If Patient ID doesn't exist, the new patient is added to patientList (new Public/Private keys are created)
* Data from webPage will be added in format <PatID, PatInfo> for user's DocID

##### Send data to miners
Data format: <DocID, <PatID, PatInfo> >
* Encrypt PatientInfo with patient's Private Key
* Sign <PatID, [PatInfo]PKpat> with doc's Private key
* Send to miner <DocID, <PatID, [PatInfo]PKpat>PKdoc>

##### Read Data
Doctor reads only his/her data (via ID) <br />
Doctors can read data only from canonical chain
* Search all values with key == DoctorID in canonical blocks
* Verify all lines and continue working where verification is true
* Decrypt <PatID, [PatInfo]PKpat> with Pat's Public Key for each PatientID
* Add to result
* Finished result is printed to web page in a form:
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
Patient's data: Public Key, Private Key, PatID <br />
Patient knows Public Keys and DocID of doctors
#### Functionalities
##### Read his/her data
Patient reads only his/her data (via ID) <br />
Patient reads data only from canonical chain
* Enter DocID into WebPage
* Search all values with key == DoctorID in canonical blocks
* Verify all lines with a specific DocID and continue working where verification is true
* Search all key-values with key == PatID
* Decrypt <PatID, [PatInfo]PKpat> with Pat's Private Key for each PatientID
* Add to result
* Finished result is printed to web page in a form:

```
Patient ID 12345 results
DocID: 119911

Timestamp = 1556758412624725900
Patient ID: 12345, Patient Data = In order to avoid heart disease, it is important to eat a healthy diet and take regular exercise.

Timestamp = 1556758399042040600
Patient ID: 12345, Patient Data = You have to be very careful how much sugar you eat when you have diabetes.
```


### API Functionalities
#### for DoctorsApp

GET /patients <br />
Web Page with information about all patients in blockchain for a specific doctor. Data shows only canonical chain without forks.
The data contains block number, timestamp, patient ID and patient information. <br />

GET /add <br />
WebPage to write PatientID & PatientInfo <br />
POST /add <br />
Add information by doctor, send to miners

#### for PatientsApp

GET /patient <br />
WebPage to write DocID <br />
POST /patient <br />
See information about a particular patient from a particular doctor

### Additional API Functionalities (Utilities)
#### for DoctorsApp

GET /start <br />
Registation functionality for doctor's application
Creation Doc's information: public, private keys
Start Doctor's heartbeat

POST /heartbeat/receive <br />
Doctor's heartbeat allow doctors to exchange information between each other, send their addresses and join the network

POST /patientlist/receive <br />
Utility function to receive patients information. Each Patient send his/her information to doctors after registration (public-private keys and ID)

GET /show <br />
Utility web page to show doctors their network with docs list and patients list

#### for PatientsApp

GET /start <br />
Registation functionality for patient's application
Creation Patient's information: public, private keys
Send information to doctors

POST /doctorlist/receive <br />
Receive public keys from doctors, add them to doctor list to have access to verify doctor's signatures

### Implementation Details

* Was added new hashmap to mpt structure to contain signatures in []byte. This decision was made to avoid []byte data type modifications while converting them to string in signature format. Thus mpt has a new structure:
```
type MerklePatriciaTrie struct {
	db map[string]Node
	kv map[string]string
	ks map[string][]byte
	root string
}
```
* Doctor's heartbeat request a special mentioned, because in this design all doctors send heartbeat to each other, so they are all in one network and can exchange patient's information, so all doctors in network will know new patients in case they work with them.
* Each patient knows doctors private keys, so they are eligible to see their information from a particular doctor.
* First miner node: 6686; first doctor node: 8813, first patient node: 9913
* Miners algorithm changed a little: miners can produce empty blocks when there is no data in dataPool, but with some data in dataPool they have to put in block at least one data block.
* For data pool implementation was made a special structure, where DB contains data <PatID, [PatInfo]PKpat> and sign contains doctor's signature for that pair.
```
type DataPool struct {
	DB   map[string]string
	Sign []byte
	Hops int
}

type ItemQueue struct {
	Items []DataPool
	Lock  sync.RWMutex
}
```

### Use Case Diagram

![alt text](https://github.com/BettyRain/blockchain_project5/blob/master/Project5UseCaseDiagram.jpeg)

### Timeline and milestones

| N%   | Milestones         | Completion Date | Actual Completion Date |
| :-------: |:-------------: | :-------------: |  :------------- |
| 1 | Add patient data (Doc's app) | April 28th | April 28th |
| 2 | View data from blocks (Doc's app) | April 28th| April 28th |
| 3 | Make interface to adding data (web page) | April 30th | May 1st |
| 4 | Checkpoint | May 1st | May 1st |
| 5 | Data exchange between applications | May 8th | May 10th |
| 6 | Data confidentiality | May 8th | May 10th |
| 7 | Data integrity | May 10th | May 14th |
| 8 | View data by patients (Patient's App) | May 10th | May 14th |

### References
Idea was found at [blockchain ideas](https://www.connectbit.com/blockchain-applications/) website.
RSA library functions were found at [blockchain documentation](https://golang.org/pkg/crypto/rsa/) and [implementation_example](https://gist.github.com/miguelmota/3ea9286bd1d3c2a985b67cac4ba2130a) websites.

### Video Link
Best Version: 