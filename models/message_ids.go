package models

var ErrorMessages = map[int]string{
	1:  "Invalid username or password.",
	2:  "The array is empty.",
	3:  "The array length exceeds 100.",
	4:  "Sender, receiver, and message arrays do not match in length.",
	5:  "Unable to retrieve new messages.",
	6:  "Account is inactive. Either the username or password is incorrect. If you recently activated the web service, please reset your web service password in the settings menu.",
	7:  "Access to the specified line is not available.",
	8:  "Invalid recipient number.",
	9:  "Insufficient account balance.",
	10: "A system error occurred. Please try again.",
	11: "Invalid IP address.",
	20: "Recipient number is filtered.",
	21: "Connection to the service provider has been lost.",
}
