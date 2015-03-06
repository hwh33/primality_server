/* The registrar is responsible for handling all user authentication.

The registrar is represented by a struct and only one registrar exists at a
time (maintained in the server). This registrar is the only one with access to
the file holding user login data.

When a user is registered, their data is stored in the registrar file as
a single line as follows:
	username,netID,passwordCheckSum
Because of this formatting, usernames and net IDs can contain no commas.

*/

package registrar

import (
	"bufio"
	"errors"
	"hash/adler32"
	"os"
	"strconv"
	"strings"
)

type userData struct {
	loggedIn         bool
	netID            string
	passwordCheckSum uint32
}

type Registrar struct {
	registrarFileName string
	users             map[string]userData
	registeredNetIDs  map[string]bool
}

/* Returns a new registrar with no initialized data. The registrar file will be
 * created with the given name. Any existing files of this name will be
 * overwritten.
 */
func NewRegistrar(registrarFileName string) (*Registrar, error) {
	registrarFile, err := os.Create(registrarFileName)
	if err != nil {
		return nil, err
	}
	registrarFile.Close()
	usersMap := make(map[string]userData)
	netIDsMap := make(map[string]bool)
	newRegistrar := Registrar{
		registrarFileName: registrarFileName,
		users:             usersMap,
		registeredNetIDs:  netIDsMap,
	}
	return &newRegistrar, nil
}

/* Reads a file with user data and uses it to create a return a registrar .*/
func NewRegistrarFromFile(registrarFileName string) (*Registrar, error) {
	usersMap := make(map[string]userData)
	netIDsMap := make(map[string]bool)

	registrarFile, err := os.Open(registrarFileName)
	if err != nil {
		return nil, err
	}
	defer registrarFile.Close()
	scanner := bufio.NewScanner(registrarFile)
	errorOccurred := false
	for scanner.Scan() {
		userDataSlice := strings.Split(scanner.Text(), ",")
		if len(userDataSlice) != 3 {
			errorOccurred = true
		} else if _, err := strconv.Atoi(userDataSlice[2]); err != nil {
			errorOccurred = true
		} else {
			username := userDataSlice[0]
			netID := userDataSlice[1]
			checkSumString, _ := strconv.Atoi((userDataSlice[2]))
			passwordCheckSum := uint32(checkSumString)
			newUserData := userData{
				loggedIn:         false,
				netID:            netID,
				passwordCheckSum: passwordCheckSum,
			}
			usersMap[username] = newUserData
			netIDsMap[netID] = true
		}
	}

	newRegistrar := Registrar{
		registrarFileName: registrarFileName,
		users:             usersMap,
		registeredNetIDs:  netIDsMap,
	}

	if errorOccurred {
		errorString := ("Registrar file possibly corrupted")
		return &newRegistrar, errors.New(errorString)
	} else {
		return &newRegistrar, nil
	}
}

/*Registers a NEW user. Should not be used for existing users. */
func (r *Registrar) RegisterUser(username, netID, password string) error {
	// Check to make sure that neither the username nor netID already exist.
	if _, exists := r.users[username]; exists {
		return errors.New("Username " + username + " has already been taken")
	} else if _, exists := r.registeredNetIDs[netID]; exists {
		return errors.New("Net ID " + netID + " is already registered")
	}

	// Screen the username and net ID for commas.
	if strings.Contains(username+netID, ",") {
		return errors.New("No commas allowed in username or net ID")
	}

	// Write the user's data to file first.
	registrarFile, err := os.OpenFile(r.registrarFileName, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer registrarFile.Close()
	passwordCheckSum := adler32.Checksum([]byte(password))
	userDataString := username + "," + netID + "," + string(passwordCheckSum) + "\n"
	_, err = registrarFile.WriteString(userDataString)
	if err != nil {
		return err
	}

	// Now push the data to the in-memory data structures and return.
	newUserData := userData{
		loggedIn:         false,
		netID:            netID,
		passwordCheckSum: passwordCheckSum,
	}
	r.users[username] = newUserData
	r.registeredNetIDs[netID] = true
	return nil
}

/* Removes a user. */
func (r *Registrar) RemoveUser(username string) error {
	// We scan the old file and copy each line to a new temporary file. When we
	// find the user data we are looking for, we simply neglect to copy it.
	tempFile, err := os.OpenFile(r.registrarFileName+"-temp", os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer tempFile.Close()
	registrarFile, err := os.Open(r.registrarFileName)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(registrarFile)
	for scanner.Scan() {
		userDataString := scanner.Text() + "\n"
		if !strings.HasPrefix(userDataString, username+",") {
			_, err = tempFile.WriteString(userDataString)
			if err != nil {
				return err
			}
		}
	}

	// Now we replace the original registrar file with the temporary file.
	registrarFile.Close()
	err = os.Remove(r.registrarFileName)
	if err != nil {
		return err
	}
	err = os.Rename(tempFile.Name(), r.registrarFileName)
	if err != nil {
		// Returning here will actually break the registrar, so we need to panic.
		panic(err.Error())
	}

	// Finally, we remove the user from the in-memory datastructures and return.
	delete(r.registeredNetIDs, r.users[username].netID)
	delete(r.users, username)
	return nil
}

/* Change a user's password. */
func (r *Registrar) ChangePassword(username, oldPW, newPW string) error {
	// First check that the username and old password are authentic.
	if !r.IsPasswordAuthentic(username, oldPW) {
		return errors.New("Invalid username and password combination")
	}

	// The easiest way to change the password is to just remove the user and
	// re-add them.
	netID := r.users[username].netID
	err := r.RemoveUser(username)
	if err != nil {
		return errors.New("A problem was encountered: " + err.Error())
	}
	err = r.RegisterUser(username, netID, newPW)
	if err != nil {
		// Returning here means that the user was removed and NOT re-added.
		return errors.New("Catastrophic problem encountered: " + err.Error() +
			". Please re-register your account.")
	}
	// At this point, we are done.
	return nil
}

/* Log the user in. */
func (r *Registrar) Login(username, password string) error {
	if !r.IsPasswordAuthentic(username, password) {
		return errors.New("Invalid username and password combination")
	}
	// We have to use a work-around to assign to a value in a map becuase of a
	// bug in Go.
	temp := r.users[username]
	temp.loggedIn = true
	r.users[username] = temp
	return nil
}

/* Log the user out. */
func (r *Registrar) Logout(username string) {
	// We have to use a work-around to assign to a value in a map becuase of a
	// bug in Go.
	temp := r.users[username]
	temp.loggedIn = false
	r.users[username] = temp
}

/* Check a password against the password registered for the user. */
func (r *Registrar) IsPasswordAuthentic(username, password string) bool {
	checkSum := adler32.Checksum([]byte(password))
	trueCheckSum := r.users[username].passwordCheckSum
	return checkSum == trueCheckSum
}
