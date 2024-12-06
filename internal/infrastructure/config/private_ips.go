package config

import (
	"net"
)

var PrivateIPs = []net.IPNet{
	{IP: net.IP{10, 0, 0, 0}, Mask: net.CIDRMask(8, 32)},     // Private range A
	{IP: net.IP{172, 16, 0, 0}, Mask: net.CIDRMask(12, 32)},  // Private range B
	{IP: net.IP{192, 168, 0, 0}, Mask: net.CIDRMask(16, 32)}, // Private range C
	{IP: net.IP{169, 254, 0, 0}, Mask: net.CIDRMask(16, 32)}, // Link-local addresses
	{IP: net.IP{224, 0, 0, 0}, Mask: net.CIDRMask(4, 32)},    // Multicast addresses
	{IP: net.IP{240, 0, 0, 0}, Mask: net.CIDRMask(4, 32)},    // Reserved addresses
	{IP: net.IP{127, 0, 0, 0}, Mask: net.CIDRMask(8, 32)},    // Loopback addresses
	{IP: net.IP{192, 0, 2, 0}, Mask: net.CIDRMask(24, 32)},   // Documentation Network
}
