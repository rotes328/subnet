package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func prepend(x []int, times int) []int {
	for i := 0; i < times; i++ {
		x = append(x, 0)
		copy(x[1:], x)
		x[0] = 0
	}
	return x
}

func convertOctetToBinary(octet uint64) []int {
	bits := []int{}
	for _, b := range strconv.FormatUint(octet, 2) {
		bits = append(bits, int(b-rune('0')))
	}
	numberOfPrepends := 8 - len(bits)
	bits = prepend(bits, numberOfPrepends)
	return bits
}

func convertMaskToBinaryList(mask []int) []int {

	octet1Bin := convertOctetToBinary(uint64(mask[0]))
	octet2Bin := convertOctetToBinary(uint64(mask[1]))
	octet3Bin := convertOctetToBinary(uint64(mask[2]))
	octet4Bin := convertOctetToBinary(uint64(mask[3]))

	maskList := []int{}
	maskList = append(octet1Bin, octet2Bin...)
	maskList = append(maskList, octet3Bin...)
	maskList = append(maskList, octet4Bin...)

	return maskList
}

func convertIPToBinaryList(ip []int) []int {

	octet1Bin := convertOctetToBinary(uint64(ip[0]))
	octet2Bin := convertOctetToBinary(uint64(ip[1]))
	octet3Bin := convertOctetToBinary(uint64(ip[2]))
	octet4Bin := convertOctetToBinary(uint64(ip[3]))

	ipList := []int{}
	ipList = append(octet1Bin, octet2Bin...)
	ipList = append(ipList, octet3Bin...)
	ipList = append(ipList, octet4Bin...)

	return ipList
}

func convertIPToBinary(ip []int) int {

	ipListString := []string{}

	// Convert to String
	for i := 0; i < 32; i++ {
		ipListString = append(ipListString, strconv.FormatInt(int64(ip[i]), 10))
	}
	ipString := strings.Join(ipListString, "")

	// Convert to Decimal for return
	fullIPBinary, err := strconv.ParseInt(ipString, 2, 0)

	if err != nil {
		invalidInput(0)
	}

	return int(fullIPBinary)

}

func validateMask(maskAsBinaryList []int) int {
	if maskAsBinaryList[0] == 0 {
		invalidInput(1)
	}
	for i := 1; i < 32; i++ {
		if maskAsBinaryList[i] > maskAsBinaryList[i-1] {
			invalidInput(1)
		}
	}
	x := 1
	for i := 1; i < 32; i++ {
		if maskAsBinaryList[i] == 1 {
			x++
		}
	}
	/* Handle Special Cases
	   Return 1 for /31
	   Return 2 for /1 - /7 (supernets)
	   Return 0 for normal range */
	if x < 8 {
		return 2
	} else if x == 32 {
		invalidInput(4)
	} else if x > 32 {
		invalidInput(1)
	} else if x == 31 {
		return 1
	}
	return 0
}

func invalidInput(errno int) {
	switch {
	case errno == 0:
		fmt.Println("Please supply an argument, --help for help.")
	case errno == 1:
		fmt.Println("Invalid mask.")
	case errno == 2:
		fmt.Println("Invalid IP address.")
	case errno == 3:
		fmt.Println("Fatal error.")
	case errno == 4:
		fmt.Println("Host address.")
	case errno == 5:
		fmt.Println("Please include IP and mask, --help for help.")
	default:
		fmt.Println("Please supply an argument, --help for help.")
	}

	os.Exit(0)
}

func checkIP(ip string) {
	if net.ParseIP(ip) == nil {
		invalidInput(2)
	}
}

func checkMaskValid(mask string) {
	if net.ParseIP(mask) == nil {
		invalidInput(1)
	}
}

func checkMask(bitmask string) int {
	mask, err := strconv.Atoi(bitmask)
	if err != nil {
		invalidInput(1)
	}
	if mask == 32 {
		invalidInput(4)
	} else if mask > 32 {
		invalidInput(1)
	}
	return mask
}

func convertMaskToBinList(bitmask int) []int {
	maskLen := bitmask
	maskBin := []int{}
	for i := 0; i < 32; i++ {
		if maskLen > 0 {
			maskBin = append(maskBin, 1)
		} else {
			maskBin = append(maskBin, 0)
		}
		maskLen--
	}
	return maskBin
}

func getSubnet(ipAsBinaryList []int, maskAsBinaryList []int) []int {
	subnetAsBinaryList := []int{}
	// Logical AND for network address
	for i := 0; i < 32; i++ {
		number := ipAsBinaryList[i] & maskAsBinaryList[i]
		subnetAsBinaryList = append(subnetAsBinaryList, number)
	}
	return subnetAsBinaryList
}

func getFirstIP(subnetDD string) string {
	x := splitStringDD(subnetDD)
	lastOctet, _ := strconv.Atoi(x[3])
	x[3] = strconv.Itoa(lastOctet + 1)
	firstIPDD := strings.Join(x, ".")
	return firstIPDD
}

func getLastIP(broadcastDD string) string {
	x := splitStringDD(broadcastDD)
	lastOctet, _ := strconv.Atoi(x[3])
	x[3] = strconv.Itoa(lastOctet - 1)
	lastIPDD := strings.Join(x, ".")
	return lastIPDD
}

func getWildcard(maskAsBinaryList []int) []int {
	wildcard := []int{}
	for i := 0; i < 32; i++ {
		if maskAsBinaryList[i] == 1 {
			wildcard = append(wildcard, 0)
		} else {
			wildcard = append(wildcard, 1)
		}
	}
	return wildcard
}

func getBroadcast(subnetAsBinaryList []int, maskAsBinaryList []int) []int {
	broadcastAsBinaryList := []int{}
	// Logical NOT on mask
	wildcard := getWildcard(maskAsBinaryList)
	for i := 0; i < 32; i++ {
		number := subnetAsBinaryList[i] + wildcard[i]
		broadcastAsBinaryList = append(broadcastAsBinaryList, number)
	}
	return broadcastAsBinaryList
}

func convertOctetListToDec(octetList []int) int {
	listString := []string{}
	for i := 0; i < 8; i++ {
		listString = append(listString, strconv.FormatInt(int64(octetList[i]), 10))
	}
	octetString := strings.Join(listString, "")
	octetAsDec, err := strconv.ParseInt(octetString, 2, 64)
	if err != nil {
		invalidInput(3)
	}
	return int(octetAsDec)
}

func getRange(slash31 bool, subnetDD string, broadcastDD string) (string, string) {
	if slash31 == true {
		return subnetDD, broadcastDD
	}
	firstIPDD := getFirstIP(subnetDD)
	lastIPDD := getLastIP(broadcastDD)
	return firstIPDD, lastIPDD
}

func convertBinaryListToDD(binaryList []int) string {
	octet1List := binaryList[0:8]
	octet2List := binaryList[8:16]
	octet3List := binaryList[16:24]
	octet4List := binaryList[24:32]

	octet1Dec := convertOctetListToDec(octet1List)
	octet2Dec := convertOctetListToDec(octet2List)
	octet3Dec := convertOctetListToDec(octet3List)
	octet4Dec := convertOctetListToDec(octet4List)

	octet1Str := strconv.Itoa(octet1Dec)
	octet2Str := strconv.Itoa(octet2Dec)
	octet3Str := strconv.Itoa(octet3Dec)
	octet4Str := strconv.Itoa(octet4Dec)

	var dd []string
	dd = append(dd, octet1Str, octet2Str, octet3Str, octet4Str)
	dottedDecimal := strings.Join(dd, ".")

	return dottedDecimal

}

func help() {
	fmt.Println("|---------------------------------------------------------|")
	fmt.Println("|                   Subnet Utility Help                   |")
	fmt.Println("|---------------------------------------------------------|")
	fmt.Println("| Supply IP and mask in slash notation or dotted decimal. |")
	fmt.Println("| Example: $subnet 10.1.1.1/24                            |")
	fmt.Println("|         $subnet 192.168.20.1 255.255.255.0              |")
	fmt.Println("|---------------------------------------------------------|")
	os.Exit(0)
}

func splitStringDD(ip string) []string {
	return strings.Split(ip, ".")
}

func convertDDtoInt(ip []string) []int {
	byteList := []int{}
	for i := 0; i < 4; i++ {
		byteInt, err := strconv.Atoi(ip[i])
		if err != nil {
			invalidInput(2)
		}
		byteList = append(byteList, byteInt)
	}
	return byteList
}

func output(ipAddress string, maskDD string, subnetDD string, broadcastDD string, firstIP string, lastIP string, supernet bool) {
	// Pad deliminator
	deliminator := "-"
	x := 29 + len(firstIP) + len(lastIP)
	padd := strings.Repeat(deliminator, x)
	y := x - 19 - len(ipAddress) - len(maskDD)
	pad1 := strings.Repeat(" ", y)

	if supernet {

		y = x - 26 - len(firstIP) - len(lastIP)
		fmt.Println(y)
		padC := strings.Repeat(" ", y)
		fmt.Printf("|%s|\n", padd)
		fmt.Printf("| For IP %s and mask %s:%s|\n", ipAddress, maskDD, pad1)
		fmt.Printf("|%s|\n", padd)
		fmt.Printf("| CIDR Range:\t\t%s - %s%s|\n", firstIP, lastIP, padC)
		fmt.Printf("|%s|\n", padd)
		return
	}

	y = x - 23 - len(subnetDD)
	pad2 := strings.Repeat(" ", y)
	y = x - 23 - len(broadcastDD)
	pad3 := strings.Repeat(" ", y)
	y = x - 26 - len(firstIP) - len(lastIP)
	pad4 := strings.Repeat(" ", y)

	fmt.Printf("|%s|\n", padd)
	fmt.Printf("| For IP %s and mask %s:%s|\n", ipAddress, maskDD, pad1)
	fmt.Printf("|%s|\n", padd)
	fmt.Printf("| Network Address:\t%s%s|\n", subnetDD, pad2)
	fmt.Printf("| Broadcast Address:\t%s%s|\n", broadcastDD, pad3)
	fmt.Printf("| Range:\t\t%s - %s%s|\n", firstIP, lastIP, pad4)
	fmt.Printf("|%s|\n", padd)
}

func handleargs(args []string) (string, []int) {
	var ipAddress string
	var mask string
	var maskAsBinaryList []int
	if len(args) < 2 {
		invalidInput(0)
	} else if len(args) == 2 {
		if strings.Contains(args[1], "/") {
			inputList := strings.Split(args[1], "/")
			ipAddress = inputList[0]
			mask = inputList[1]
		} else if args[1] == "--help" {
			help()
		} else {
			invalidInput(5)
		}
	} else {
		ipAddress = args[1]
		mask = args[2]
	}
	if len(mask) < 3 {
		mask := checkMask(mask)
		maskAsBinaryList = convertMaskToBinList(mask)
	} else {
		checkMaskValid(mask)
		maskStr := splitStringDD(mask)
		maskInt := convertDDtoInt(maskStr)
		maskAsBinaryList = convertMaskToBinaryList(maskInt)
	}
	return ipAddress, maskAsBinaryList
}

func getIPAsBinaryList(ipAddress string) []int {
	ipAddressStr := splitStringDD(ipAddress)
	ipInt := convertDDtoInt(ipAddressStr)
	ipAsBinaryList := convertIPToBinaryList(ipInt)
	return ipAsBinaryList
}

func handleSpecialMasks(maskType int) (bool, bool) {
	switch {
	case maskType == 1: // /31 Mask
		return true, false
	case maskType == 2: // /1 - /7 Mask
		return false, true
	}
	return false, false // No special cases
}

func subnetCalc(ipAddress string, maskAsBinaryList []int) (string, string, string, string, string, bool) {
	maskType := validateMask(maskAsBinaryList)
	slash31, supernet := handleSpecialMasks(maskType)
	ipAsBinaryList := getIPAsBinaryList(ipAddress)
	subnetAsBinaryList := getSubnet(ipAsBinaryList, maskAsBinaryList)
	broadcastAsBinaryList := getBroadcast(subnetAsBinaryList, maskAsBinaryList)
	subnetDD := convertBinaryListToDD(subnetAsBinaryList)
	broadcastDD := convertBinaryListToDD(broadcastAsBinaryList)
	maskDD := convertBinaryListToDD(maskAsBinaryList)
	firstIP, lastIP := getRange(slash31, subnetDD, broadcastDD)
	return maskDD, subnetDD, broadcastDD, firstIP, lastIP, supernet
}

func main() {

	var ipAddress string
	var maskAsBinaryList []int

	// Handle arguments
	args := os.Args
	ipAddress, maskAsBinaryList = handleargs(args)

	// Validate IP Address
	checkIP(ipAddress)

	// Do bitmath
	maskDD, subnetDD, broadcastDD, firstIP, lastIP, supernet := subnetCalc(ipAddress, maskAsBinaryList)

	// Render output
	output(ipAddress, maskDD, subnetDD, broadcastDD, firstIP, lastIP, supernet)

}
