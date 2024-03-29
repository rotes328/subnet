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

func handleMulticast(class int) bool {
	if class == 3 { // Class D
		return true
	} else if class > 3 { // Class E
		invalidInput(2)
	}
	return false
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

func validateMask(maskAsBinaryList []int) {
	if maskAsBinaryList[0] == 0 {
		invalidInput(1)
	}
	for i := 1; i < 32; i++ {
		if maskAsBinaryList[i] > maskAsBinaryList[i-1] {
			invalidInput(1)
		}
	}
}

func getMaskClass(maskAsBinaryList []int) int {
	/* Handle Special Cases and identify class boundary
	    Return 0 for class A valid
		Return 1 for class B valid
		Return 2 for class C valid
		Return 3 for /4 - /7 (can be valid for multicast)
		Return 4 for /1 - /3 (always supernet)
		Return 5 for /31
		Return 6 for /32
	*/

	x := 1
	for i := 1; i < 32; i++ {
		if maskAsBinaryList[i] == 1 {
			x++
		}
	}

	if x > 32 {
		invalidInput(1)
	} else if x == 32 {
		return 6
	} else if x == 31 {
		return 5
	} else if (x >= 24) && (x <= 30) {
		return 2
	} else if (x >= 16) && (x <= 23) {
		return 1
	} else if (x >= 8) && (x <= 15) {
		return 0
	} else if (x >= 4) && (x <= 7) {
		return 3
	}
	return 4
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
	case errno == 6:
		fmt.Println("Multicast IP, not supported.")
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
	if mask > 32 {
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

func getRange(subnetDD string, broadcastDD string, slash31 bool, multicast bool, supernet bool) (string, string) {
	if slash31 || multicast || supernet {
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
	fmt.Println("╔═════════════════════════════════════════════════════════╗")
	fmt.Println("║                   Subnet Utility Help                   ║")
	fmt.Println("╟─────────────────────────────────────────────────────────╢")
	fmt.Println("║ Supply IP and mask in slash notation or dotted decimal. ║")
	fmt.Println("║ Example: $subnet 10.1.1.1/24                            ║")
	fmt.Println("║         $subnet 192.168.20.1 255.255.255.0              ║")
	fmt.Println("╚═════════════════════════════════════════════════════════╝")
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

func output(ipAddress string, maskDD string, subnetDD string, broadcastDD string, firstIP string, lastIP string, supernet bool, multicast bool) {
	// Pad deliminator
	outsideDeliminator := "═"
	deliminator := "─"
	x := 29 + len(firstIP) + len(lastIP)
	padD := strings.Repeat(deliminator, x)
	padO := strings.Repeat(outsideDeliminator, x)
	y := x - 19 - len(ipAddress) - len(maskDD)
	pad1 := strings.Repeat(" ", y)

	if supernet {

		y = x - 26 - len(firstIP) - len(lastIP)
		padC := strings.Repeat(" ", y)
		fmt.Printf("╔%s╗\n", padO)
		fmt.Printf("║ For IP %s and mask %s:%s║\n", ipAddress, maskDD, pad1)
		fmt.Printf("╟%s╢\n", padD)
		fmt.Printf("║ CIDR Range:\t\t%s - %s%s║\n", firstIP, lastIP, padC)
		fmt.Printf("╚%s╝\n", padO)
		return
	}

	if multicast {

		y = x - 26 - len(firstIP) - len(lastIP)
		padC := strings.Repeat(" ", y)
		fmt.Printf("╔%s╗\n", padO)
		fmt.Printf("║ For IP %s and mask %s:%s║\n", ipAddress, maskDD, pad1)
		fmt.Printf("╟%s╢\n", padD)
		fmt.Printf("║ Multicast Range:\t%s - %s%s║\n", firstIP, lastIP, padC)
		fmt.Printf("╚%s╝\n", padO)
		return
	}

	y = x - 23 - len(subnetDD)
	pad2 := strings.Repeat(" ", y)
	y = x - 23 - len(broadcastDD)
	pad3 := strings.Repeat(" ", y)
	y = x - 26 - len(firstIP) - len(lastIP)
	pad4 := strings.Repeat(" ", y)

	fmt.Printf("╔%s╗\n", padO)
	fmt.Printf("║ For IP %s and mask %s:%s║\n", ipAddress, maskDD, pad1)
	fmt.Printf("╟%s╢\n", padD)
	fmt.Printf("║ Network Address:\t%s%s║\n", subnetDD, pad2)
	fmt.Printf("║ Broadcast Address:\t%s%s║\n", broadcastDD, pad3)
	fmt.Printf("║ Range:\t\t%s - %s%s║\n", firstIP, lastIP, pad4)
	fmt.Printf("╚%s╝\n", padO)
}

func checkClass(ipAsBinaryList []int) int {
	/* Return 0 for class A
	   Return 1 for class B
	   Return 2 for class C
	   Return 3 for multicast
	   Return 4 for invalid */
	for i := 0; i < 4; i++ {
		if ipAsBinaryList[i] == 0 {
			return i
		}
	}
	return 4
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

func handleSlash31(maskClass int) bool {
	// Return true for /31 mask
	return maskClass == 5
}

func determineSupernet(ipClass int, maskClass int) bool {
	// Returns true for supernet

	if (ipClass == 3) && (maskClass < 4) {
		return false // Multicast range
	} else if maskClass == 5 {
		return false // /31 Range
	} else if maskClass >= 3 {
		return true // Always supernet if not multicast
	} else if ipClass > maskClass {
		return true // Supernet if true
	}
	return false
}

func hostOutput(ipAddress string) {
	padO := strings.Repeat("═", len(ipAddress))
	fmt.Printf("╔═════════════%s════╗\n", padO)
	fmt.Printf("║ IP Address: %s/32 ║\n", ipAddress)
	fmt.Printf("╚═════════════%s════╝\n", padO)
	os.Exit(0)
}

func subnetCalc(ipAddress string, maskAsBinaryList []int, maskClass int) (string, string, string, string, string, bool, bool) {
	validateMask(maskAsBinaryList)
	ipAsBinaryList := getIPAsBinaryList(ipAddress)
	ipClass := checkClass(ipAsBinaryList)
	multicast := handleMulticast(ipClass)
	slash31 := handleSlash31(maskClass)
	supernet := determineSupernet(ipClass, maskClass)
	subnetAsBinaryList := getSubnet(ipAsBinaryList, maskAsBinaryList)
	broadcastAsBinaryList := getBroadcast(subnetAsBinaryList, maskAsBinaryList)
	subnetDD := convertBinaryListToDD(subnetAsBinaryList)
	broadcastDD := convertBinaryListToDD(broadcastAsBinaryList)
	maskDD := convertBinaryListToDD(maskAsBinaryList)
	firstIP, lastIP := getRange(subnetDD, broadcastDD, multicast, slash31, supernet)
	return maskDD, subnetDD, broadcastDD, firstIP, lastIP, supernet, multicast
}

func main() {

	var ipAddress string
	var maskAsBinaryList []int

	// Handle arguments
	args := os.Args
	ipAddress, maskAsBinaryList = handleargs(args)

	// Validate IP Address
	checkIP(ipAddress)

	// Handle /32
	maskClass := getMaskClass(maskAsBinaryList)
	if maskClass == 6 {
		hostOutput(ipAddress)
	}

	// Do bitmath
	maskDD, subnetDD, broadcastDD, firstIP, lastIP, supernet, multicast := subnetCalc(ipAddress, maskAsBinaryList, maskClass)

	// Render output
	output(ipAddress, maskDD, subnetDD, broadcastDD, firstIP, lastIP, supernet, multicast)

}
