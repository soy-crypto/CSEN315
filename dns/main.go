package main

import (
	"bufio"
	"crypto/rand"
	"flag"
	"fmt"
	"math/big"
	"net"
	"os"
	"strings"

	"golang.org/x/net/dns/dnsmessage"
)

type RecordType uint16

const (
	TYPE_A     RecordType = 1
	TYPE_NS    RecordType = 2
	TYPE_CNAME RecordType = 5
	TYPE_TXT   RecordType = 16
	TYPE_AAAA  RecordType = 28
)

var RecordTypes map[string]RecordType = map[string]RecordType{
	"A":     TYPE_A,
	"NS":    TYPE_NS,
	"CNAME": TYPE_CNAME,
	"TXT":   TYPE_TXT,
	"AAAA":  TYPE_AAAA,
}

func main() {
	// get all command line arguments
	names := os.Args[1:]
	t := flag.String("t", "A", "the record type to query for each name")
	flag.Parse()

	// input validation
	if len(names) == 0 {
		fmt.Println("Not enough arguments, must pass in at least one name")
		os.Exit(1)
	}

	if _, exists := RecordTypes[*t]; !exists {
		keys := make([]string, 0, len(RecordTypes))
		for k := range RecordTypes {
			keys = append(keys, k)
		}
		fmt.Printf("Specified record type %s doesn't exist. Must be one of %v", *t, keys)
		os.Exit(1)
	}

	// Invoke the resolve function for each of the given names
	for _, name := range names {
		fmt.Printf("%s,%s\n", name, strings.Join(resolve(name, RecordTypes[*t]), ""))
	}

	fmt.Printf("\n")
}

// Resolver
func resolve(name string, t RecordType) []string {
	// most of your code should go here. use a switch statement
	// so each resolution type goes into a different function
	resolvedValue := make([]string, 0, 100)
	hostName := name + "."

	//Resolve the name
	/* check what kind of data to be requested to name serverse*/
	var typ dnsmessage.Type
	switch t {
	//enquire about ipv4 address
	case TYPE_A:
		typ = 1

	//enquire about ipv6 address
	case TYPE_AAAA:
		typ = 28

	//Enquire about
	case TYPE_CNAME:
		typ = 5

	case TYPE_NS:
		typ = 2

	case TYPE_TXT:
		typ = 16

	default:
		fmt.Printf("Unsupported record type: %v\n", t)
	}

	/* Do query */
	address, err := query(hostName, typ)
	if err != nil {
		fmt.Printf("Error resolving A record for %s: %v\n", name, err)
		return resolvedValue
	}

	resolvedValue = append(resolvedValue, address) // Convert IP to string

	//Result
	var result []string = make([]string, 1)
	result[0] = resolvedValue[0]

	//Return
	return result

}

// ************************************** Newly Added ***************************************************
// all the address of root servers
const ROOT_SERVERS = "198.41.0.4,199.9.14.201,192.33.4.12,199.7.91.13,192.203.230.10,192.5.5.241,192.112.36.4,198.97.190.53"

// convert ROOT_SERVERS to an array of root servers
func getRootServers() []net.IP {
	rootServers := []net.IP{}
	for _, rootServer := range strings.Split(ROOT_SERVERS, ",") {
		rootServers = append(rootServers, net.ParseIP(rootServer))
	}
	return rootServers
}

// Do query of all the root servers
func query(name string, TYPE dnsmessage.Type) (string, error) {

	//Do dns query
	/* build a question*/
	Question := dnsmessage.Question{
		Name:  dnsmessage.MustNewName(name),
		Type:  TYPE,
		Class: dnsmessage.ClassINET,
	}

	/* call the function dnsQuery to get a response from all name servers including root server*/
	response, err := dnsQuery(getRootServers(), Question)
	if err != nil {
		fmt.Println("rand error: %s", err)
	}

	var INFO string
	for _, answer := range response.Answers {
		if TYPE == 1 {
			if strings.Contains(answer.Body.GoString(), "{A: [") {
				originalStr := strings.Split(strings.Split(answer.Body.GoString(), "}")[0], "{")[2]
				ip := strings.Replace(originalStr, ", ", ".", -1)
				INFO = ip
				break
			}

		} else if TYPE == 28 {
			if strings.Contains(answer.Body.GoString(), "{AAAA: [") {
				originalStr := strings.Split(strings.Split(answer.Body.GoString(), "}")[0], "{")[2]
				ip := strings.Replace(originalStr, ", ", ".", -1)
				//fmt.Printf("IP %s\n", ip)
				INFO = ip
				break
			}

		} else if TYPE == 2 {
			if strings.Contains(answer.Body.GoString(), "NS: ") {
				originalStr := strings.Split(strings.Split(strings.Split(answer.Body.GoString(), "}")[0], "(")[1], ")")[0]
				ip := strings.Replace(originalStr, ", ", ".", -1)
				INFO = ip

				/* resovle the name sever */
				/* Do query */
				INFO = INFO[1:len(INFO)-2] + "."
				address, err := query(INFO, 1)
				if err != nil {
					fmt.Printf("Error resolving A record for %s: %v\n", name, err)
					return INFO, nil
				}

				INFO += ", " + address
				break
			}

		} else if TYPE == 5 {
			if strings.Contains(answer.Body.GoString(), "CNAME: ") {
				originalStr := strings.Split(strings.Split(strings.Split(answer.Body.GoString(), "}")[0], "(")[1], ")")[0]
				ip := strings.Replace(originalStr, ", ", ".", -1)
				INFO = ip

				/* resovle the name sever */
				/* Do query */
				INFO = INFO[1:len(INFO)-2] + "."
				address, err := query(INFO, 1)
				if err != nil {
					fmt.Printf("Error resolving A record for %s: %v\n", name, err)
					return INFO, nil
				}

				INFO += ", " + address
				break
			}
		} else if TYPE == 16 {
			if strings.Contains(answer.Body.GoString(), "TXT: ") {
				//fmt.Print("TXT:", answer.Body.GoString())
				originalStr := strings.Split(strings.Split(strings.Split(answer.Body.GoString(), "}")[0], "{")[2], "\\")[1]
				ip := strings.Replace(originalStr, ", ", ".", -1)
				INFO = ip
				break
			}
		} //

	} //for

	//Return
	return INFO, nil

} //func

func dnsQuery(servers []net.IP, question dnsmessage.Question) (*dnsmessage.Message, error) {
	//fmt.Printf("Question: %+v\n", question)
	for i := 0; i < 3; i++ {
		//call outgoingDnsQuery
		dnsAnswer, header, err := outgoingDnsQuery(servers, question)
		if err != nil {
			return nil, err
		}

		/* Get retunr ansers */
		parsedAnswers, err := dnsAnswer.AllAnswers()
		if err != nil {
			return nil, err
		}

		if header.Authoritative {
			return &dnsmessage.Message{
				Header:  dnsmessage.Header{Response: true},
				Answers: parsedAnswers,
			}, nil

		}

		/* Get retunr authorities */
		authorities, err := dnsAnswer.AllAuthorities()
		if err != nil {
			return nil, err
		}

		if len(authorities) == 0 {
			return &dnsmessage.Message{
				Header: dnsmessage.Header{RCode: dnsmessage.RCodeNameError},
			}, nil
		}

		/* Get all the nameserveres*/
		nameservers := make([]string, len(authorities))
		for k, authority := range authorities {
			if authority.Header.Type == dnsmessage.TypeNS {
				//fmt.Printf("authority.Body.(*dnsmessage.NSResource).NS: %s \n", authority.Body.(*dnsmessage.NSResource).NS)
				nameservers[k] = authority.Body.(*dnsmessage.NSResource).NS.String()
				//fmt.Printf("nameserver: %s \n", nameservers[k])
			}
		}

		/* Get all the additional coresponding to all the authorities */
		additionals, err := dnsAnswer.AllAdditionals()
		if err != nil {
			return nil, err
		}
		newResolverServersFound := false
		servers = []net.IP{} // set servers as empty
		for _, additional := range additionals {
			if additional.Header.Type == dnsmessage.TypeA {
				for _, nameserver := range nameservers {
					if additional.Header.Name.String() == nameserver {
						newResolverServersFound = true
						servers = append(servers, additional.Body.(*dnsmessage.AResource).A[:])
					} //if

				} //for

			} //if

		} //for

		if !newResolverServersFound {
			for _, nameserver := range nameservers {
				if !newResolverServersFound {
					response, err := dnsQuery(getRootServers(), dnsmessage.Question{Name: dnsmessage.MustNewName(nameserver), Type: dnsmessage.TypeA, Class: dnsmessage.ClassINET})
					if err != nil {
						fmt.Printf("warning: lookup of nameserver %s failed: %err\n", nameserver, err)
					} else {
						newResolverServersFound = true
						for _, answer := range response.Answers {

							if answer.Header.Type == dnsmessage.TypeA {
								servers = append(servers, answer.Body.(*dnsmessage.AResource).A[:])
							}

						} //for

					} //else

				} //if

			} //for

		} //if

		//fmt.Printf("%s", newResolverServersFound)

	}

	return &dnsmessage.Message{
		Header: dnsmessage.Header{RCode: dnsmessage.RCodeServerFailure},
	}, nil

}

func outgoingDnsQuery(servers []net.IP, question dnsmessage.Question) (*dnsmessage.Parser, *dnsmessage.Header, error) {
	/*used for randomly choosing a random number*/
	max := ^uint16(0)
	randomNumber, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return nil, nil, err
	}

	/* build a new message */
	message := dnsmessage.Message{
		Header: dnsmessage.Header{
			ID:       uint16(randomNumber.Int64()),
			Response: false,
			OpCode:   dnsmessage.OpCode(0),
		},
		Questions: []dnsmessage.Question{question},
	}

	/* encode the new message */
	buf, err := message.Pack()
	if err != nil {
		return nil, nil, err
	}

	/* Find one root server avaiable now */
	var conn net.Conn
	for _, server := range servers {
		conn, err = net.Dial("udp", server.String()+":53")
		if err == nil {
			break
		}
	}
	if conn == nil {
		return nil, nil, fmt.Errorf("failed to make connection to servers: %s", err)
	}

	/* send the new message to the choosen server */
	_, err = conn.Write(buf)
	if err != nil {
		return nil, nil, err
	}

	/* receive the answer to the message from the choosen server */
	answer := make([]byte, 512)
	n, err := bufio.NewReader(conn).Read(answer)
	if err != nil {
		return nil, nil, err
	}

	/* Close the connection */
	conn.Close()

	//parse the answer
	var p dnsmessage.Parser
	/* Get the head part of the answer */
	header, err := p.Start(answer[:n])
	if err != nil {
		return nil, nil, fmt.Errorf("parser start error: %s", err)
	}

	/* Get the question part of the answer */
	questions, err := p.AllQuestions()
	if err != nil {
		return nil, nil, err
	}

	if len(questions) != len(message.Questions) {
		return nil, nil, fmt.Errorf("answer packet doesn't have the same amount of questions")
	}

	return &p, &header, nil

}