package main

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"math/big"
	"net"
	"strings"

	"golang.org/x/net/dns/dnsmessage"
)

const ROOT_SERVERS = "198.41.0.4,199.9.14.201,192.33.4.12,199.7.91.13,192.203.230.10,192.5.5.241,192.112.36.4,198.97.190.53"

func main() {
	name := "www.youtube.com."
	query(name, dnsmessage.TypeAAAA)
}

func getRootServers() []net.IP {
	rootServers := []net.IP{}
	for _, rootServer := range strings.Split(ROOT_SERVERS, ",") {
		rootServers = append(rootServers, net.ParseIP(rootServer))
	}
	return rootServers
}

func query(name string, TYPE dnsmessage.Type) {

	/* Do dns query*/
	Question := dnsmessage.Question{
		Name:  dnsmessage.MustNewName(name),
		Type:  TYPE,
		Class: dnsmessage.ClassINET,
	}

	response, err := dnsQuery(getRootServers(), Question)
	if err != nil {
		fmt.Println("rand error: %s", err)
	}

	//fmt.Println("responseBuffer: %s", response)
	fmt.Printf("@@@@@@ HEADER %+v \n", response.Header)
	for _, answer := range response.Answers {
		fmt.Printf("STR %+v %+v \n", answer.Header, answer.Body)
		fmt.Printf("%s\n", answer.Body.GoString())
		if TYPE == 1 {
			if strings.Contains(answer.Body.GoString(), "{A: [") {
				fmt.Printf("1 @@@@@@BODY  %s \n", answer.Body.GoString())
				originalStr := strings.Split(strings.Split(answer.Body.GoString(), "}")[0], "{")[2]
				ip := strings.Replace(originalStr, ", ", ".", -1)
				fmt.Printf("IP %s\n", ip)
				break
			}

		} else if TYPE == 28 {
			if strings.Contains(answer.Body.GoString(), "{AAAA: [") {
				fmt.Printf("28 @@@@@@BODY  %s \n", answer.Body.GoString())
				originalStr := strings.Split(strings.Split(answer.Body.GoString(), "}")[0], "{")[2]
				ip := strings.Replace(originalStr, ", ", ".", -1)
				fmt.Printf("IP %s\n", ip)
				break
			}

		} else if TYPE == 2 {
			if strings.Contains(answer.Body.GoString(), "NS") {
				fmt.Printf("2 @@@@@@BODY  %s \n", answer.Body.GoString())
				break
			}

		} else if TYPE == 5 {
			if strings.Contains(answer.Body.GoString(), "CNAME") {
				fmt.Printf("5 @@@@@@BODY  %s \n", answer.Body.GoString())
				break
			}
		} else if TYPE == 16 {
			if strings.Contains(answer.Body.GoString(), "TXT") {
				fmt.Printf("16 @@@@@@BODY  %s \n", answer.Body.GoString())
				break
			}
		} //

	} //for

} //func

func dnsQuery(servers []net.IP, question dnsmessage.Question) (*dnsmessage.Message, error) {
	fmt.Printf("Question: %+v\n", question)
	for i := 0; i < 3; i++ {
		//call outgoingDnsQuery
		dnsAnswer, header, err := outgoingDnsQuery(servers, question)
		if err != nil {
			return nil, err
		}

		/* get retunr ansers */
		parsedAnswers, err := dnsAnswer.AllAnswers()
		if err != nil {
			return nil, err
		}

		if header.Authoritative {
			for _, resource := range parsedAnswers {
				fmt.Printf("!!!!!!!!!!!Found Header: %+v Body:%+v \n", resource.Header, resource.Body)
			}

			return &dnsmessage.Message{
				Header:  dnsmessage.Header{Response: true},
				Answers: parsedAnswers,
			}, nil

		}

		/* get retunr authorities */
		authorities, err := dnsAnswer.AllAuthorities()
		if err != nil {
			return nil, err
		}

		if len(authorities) == 0 {
			return &dnsmessage.Message{
				Header: dnsmessage.Header{RCode: dnsmessage.RCodeNameError},
			}, nil
		}

		/* get all the nameserveres*/
		nameservers := make([]string, len(authorities))
		for k, authority := range authorities {
			fmt.Printf("authority: %s \n", authority.Body.GoString())
			if authority.Header.Type == dnsmessage.TypeNS {
				fmt.Printf("authority.Body.(*dnsmessage.NSResource).NS: %s \n", authority.Body.(*dnsmessage.NSResource).NS)
				nameservers[k] = authority.Body.(*dnsmessage.NSResource).NS.String()
				fmt.Printf("nameserver: %s \n", nameservers[k])
			}
		}

		/* get all the additional coresponding to all the authorities */
		additionals, err := dnsAnswer.AllAdditionals()
		if err != nil {
			return nil, err
		}
		newResolverServersFound := false
		servers = []net.IP{} // set servers as empty
		for _, additional := range additionals {
			fmt.Printf("additional: %+v \n", additional)
			fmt.Printf("additional body: %s \n", additional.Body.GoString())
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
	fmt.Printf("New outgoing dns query for %s, servers: %+v\n", question.Name.String(), servers)

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

	err = p.SkipAllQuestions()
	if err != nil {
		return nil, nil, err
	}

	return &p, &header, nil

}
