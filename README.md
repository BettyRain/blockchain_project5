Project 5
================================================

### What

Healthcare data

### Why/How

Today, doctors, nurses & health professionals are limited in the level of care they can provide, because they canâ€™t view complete, accurate health records. Doctors will be able to publish medical records safely on the blockchain and only authorized people will be able to access the data.
Medical institutions will be able to securely share data, no medical records will be lost. Test results will come on time and the probability of correct diagnosis and effective treatment will be higher. Blockchain is immutable, so nobody will be able to fake or lose data. Blockchain is decentralized, so the information will be stored not only in one datacenter, in this case - in different medical institutions.

### Functionalities/Success

Doctor can add information about patient
Information will be available to view
Information will be non changeable
Patient will be able to see itâ€™s data by entering personal code

### Timeline and milestones

| N%   | Milestones         | Completion Date | Implementation |
| :-------: |:-------------: | :-------------:| :-----|
| 1 | Add patient data | April 28nd | Hashmap with key-value pair will be inserted in mpt by each pair | 
| 2 | View data from blocks (by doctor) | April 28th | Doctor can see data from all blocks (only canonical chain, no forks|
| 3 | Miners can add block with data only | April 30th | Miners should wait till they have a new data and only then solve puzzle|
| 4 | View data by personal code| May 8th | |
| 5 | Make added data immutable | May 10th | Miners can't change entered data |
| 6 | Make interface to adding data (web page) | May 10th | |


### API Functionalities

/patients
Web Page with information about all patients in blockchain. Data shows only canonical chain without forks.
The data contains block number, timestamp, patient ID and patient information.

/add
Add information by doctor

/patient
See information about a particular patient