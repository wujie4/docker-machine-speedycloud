package servers

import (
	//"encoding/base64"
	//"errors"
	//"fmt"

    "github.com/hna/speedycloud"
    //"github.com/hna/speedycloud/pagination"
    //"golang.org/x/tools/container/intsets"
    "net/url"
    "bytes"
    "fmt"
)

// List makes a request against the API to list servers accessible to you.
//func List(client *speedycloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
//	url := listDetailURL(client)
//
//	if opts != nil {
//		query, err := opts.ToServerListQuery()
//		if err != nil {
//			return pagination.Pager{Err: err}
//		}
//		url += query
//	}
//
//	createPageFn := func(r pagination.PageResult) pagination.Page {
//		return ServerPage{pagination.LinkedPageBase{PageResult: r}}
//	}
//
//	return pagination.NewPager(client, url, createPageFn)
//}

// CreateOptsBuilder describes struct types that can be accepted by the Create call.
// The CreateOpts struct in this package does.
//type CreateOptsBuilder interface {
//	ToServerCreateMap() (map[string]interface{}, error)
//}

// Network is used within CreateOpts to control a new server's network attachments.
//type Network struct {
//	// UUID of a nova-network to attach to the newly provisioned server.
//	// Required unless Port is provided.
//	UUID string
//
//	// Port of a neutron network to attach to the newly provisioned server.
//	// Required unless UUID is provided.
//	Port string
//
//	// FixedIP [optional] specifies a fixed IPv4 address to be used on this network.
//	FixedIP string
//}

// CreateOpts specifies server creation parameters.
type CreateOpts struct {
	// Name [required] is the name to assign to the newly launched server.
	Name string

	// ImageRef [required] is the ID or full URL to the image that contains the server's OS and initial state.
	// Optional if using the boot-from-volume extension.
	ImageName string

	// UserData [optional] contains configuration information or scripts to use upon launch.
	// Create will base64-encode it for you.
	//UserData []byte

	// AvailabilityZone in which to launch the server.
	AvailabilityZone string

	// Networks [optional] dictates how this server will be attached to available networks.
	// By default, the server will be attached to all isolated networks for the tenant.
	Network string

	// AdminPass [optional] sets the root user password. If not set, a randomly-generated
	// password will be created and returned in the response.
	AdminPass string
	CpuNumber int
    Memory    int
    DiskType  string
    DiskCapacity  int
    Isp       string
    Bandwidth int
    SshKey    string
    BootScript string
}

// ToServerCreateMap assembles a request body based on the contents of a CreateOpts.
func (opts CreateOpts) ToServerCreateUrlEncode() (*bytes.Buffer, error) {
	server := url.Values{}

    server.Set("az", opts.AvailabilityZone)
    //server.Set("name", opts.Name)

	server.Set("image", opts.ImageName)
	server.Set("cpu", fmt.Sprintf("%d", opts.CpuNumber))
    server.Set("memory", fmt.Sprintf("%d", opts.Memory))
    server.Set("disk_type", opts.DiskType)
    server.Set("disk", fmt.Sprintf("%d", opts.DiskCapacity))
    server.Set("isp", opts.Isp)
    server.Set("bandwidth", fmt.Sprintf("%d", opts.Bandwidth))
    server.Set("sshkeys", opts.SshKey)


	//if opts.UserData != nil {
	//	encoded := base64.StdEncoding.EncodeToString(opts.UserData)
	//	server["user_data"] = &encoded
	//}
	//if opts.Metadata != nil {
	//	server["metadata"] = opts.Metadata
	//}
	if opts.AdminPass != "" {
		server.Set("adminPass", opts.AdminPass)
	}

	if opts.Network != "" {
		server.Set("networks", opts.Network)
	}
	if opts.BootScript != ""{
        server.Set("bootscript", opts.BootScript)
    }

	return  bytes.NewBufferString(server.Encode()), nil
}

// Create requests a server to be provisioned to the user in the current tenant.
func Create(client *speedycloud.ServiceClient, opts CreateOpts) CreateResult {
	var res CreateResult

	reqBody, err := opts.ToServerCreateUrlEncode()
	if err != nil {
		res.Err = err
		return res
	}

	_, res.Err = client.Post(createURL(client), reqBody, &res.Body, nil)
	return res
}

// Delete requests that a server previously provisioned be removed from your account.
func Delete(client *speedycloud.ServiceClient, id string) DeleteResult {
	var res DeleteResult
	_, res.Err = client.Post(actionURL(client, id, "destroy"),
        bytes.NewBufferString(""),
        &res.Body,
        nil)
	return res
}

// Get requests details on a single server, by ID.
func Get(client *speedycloud.ServiceClient, id string) GetResult {
	var result GetResult
	_, result.Err = client.Post(getURL(client, id),
        bytes.NewBufferString(""),
        &result.Body,
        nil)
	return result
}

func Alias(client *speedycloud.ServiceClient, id string, aliasName string) UpdateResult {
    var result UpdateResult
    _, result.Err = client.Post(actionURL(client, id, "alias"),
        bytes.NewBufferString(fmt.Sprintf("alias=%s", aliasName)),
        &result.Body,
        nil)
    return result
}

func Group(client *speedycloud.ServiceClient, id string, groupName string) UpdateResult {
    var result UpdateResult
    _, result.Err = client.Post(actionURL(client, id, "group"),
        bytes.NewBufferString(fmt.Sprintf("group=%s", groupName)),
        &result.Body,
        nil)
    return result
}

// UpdateOptsBuilder allows extensions to add additional attributes to the Update request.
//type UpdateOptsBuilder interface {
//	ToServerUpdateMap() map[string]interface{}
//}

//// UpdateOpts specifies the base attributes that may be updated on an existing server.
//type UpdateOpts struct {
//	// Name [optional] changes the displayed name of the server.
//	// The server host name will *not* change.
//	// Server names are not constrained to be unique, even within the same tenant.
//	Name string
//
//	// AccessIPv4 [optional] provides a new IPv4 address for the instance.
//	AccessIPv4 string
//
//	// AccessIPv6 [optional] provides a new IPv6 address for the instance.
//	AccessIPv6 string
//}
//
//// ToServerUpdateMap formats an UpdateOpts structure into a request body.
//func (opts UpdateOpts) ToServerUpdateMap() map[string]interface{} {
//	server := make(map[string]string)
//	if opts.Name != "" {
//		server["name"] = opts.Name
//	}
//	if opts.AccessIPv4 != "" {
//		server["accessIPv4"] = opts.AccessIPv4
//	}
//	if opts.AccessIPv6 != "" {
//		server["accessIPv6"] = opts.AccessIPv6
//	}
//	return map[string]interface{}{"server": server}
//}
//
//// Update requests that various attributes of the indicated server be changed.
//func Update(client *speedycloud.ServiceClient, id string, opts UpdateOptsBuilder) UpdateResult {
//	var result UpdateResult
//	reqBody := opts.ToServerUpdateMap()
//	_, result.Err = client.Put(updateURL(client, id), reqBody, &result.Body, &speedycloud.RequestOpts{
//		OkCodes: []int{200},
//	})
//	return result
//}

// ChangeAdminPassword alters the administrator or root password for a specified server.
//func ChangeAdminPassword(client *speedycloud.ServiceClient, id, newPassword string) ActionResult {
//	var req struct {
//		ChangePassword struct {
//			AdminPass string `json:"adminPass"`
//		} `json:"changePassword"`
//	}
//
//	req.ChangePassword.AdminPass = newPassword
//
//	var res ActionResult
//	_, res.Err = client.Post(actionURL(client, id), req, nil, nil)
//	return res
//}

// ErrArgument errors occur when an argument supplied to a package function
// fails to fall within acceptable values.  For example, the Reboot() function
// expects the "how" parameter to be one of HardReboot or SoftReboot.  These
// constants are (currently) strings, leading someone to wonder if they can pass
// other string values instead, perhaps in an effort to break the API of their
// provider.  Reboot() returns this error in this situation.
//
// Function identifies which function was called/which function is generating
// the error.
// Argument identifies which formal argument was responsible for producing the
// error.
// Value provides the value as it was passed into the function.
//type ErrArgument struct {
//	Function, Argument string
//	Value              interface{}
//}
//
//// Error yields a useful diagnostic for debugging purposes.
//func (e *ErrArgument) Error() string {
//	return fmt.Sprintf("Bad argument in call to %s, formal parameter %s, value %#v", e.Function, e.Argument, e.Value)
//}
//
//func (e *ErrArgument) String() string {
//	return e.Error()
//}

//// RebootMethod describes the mechanisms by which a server reboot can be requested.
//type RebootMethod string
//
//// These constants determine how a server should be rebooted.
//// See the Reboot() function for further details.
//const (
//	SoftReboot RebootMethod = "SOFT"
//	HardReboot RebootMethod = "HARD"
//	OSReboot                = SoftReboot
//	PowerCycle              = HardReboot
//)

// Reboot requests that a given server reboot.
// Two methods exist for rebooting a server:
//
// HardReboot (aka PowerCycle) restarts the server instance by physically cutting power to the machine, or if a VM,
// terminating it at the hypervisor level.
// It's done. Caput. Full stop.
// Then, after a brief while, power is restored or the VM instance restarted.
//
// SoftReboot (aka OSReboot) simply tells the OS to restart under its own procedures.
// E.g., in Linux, asking it to enter runlevel 6, or executing "sudo shutdown -r now", or by asking Windows to restart the machine.
func Reboot(client *speedycloud.ServiceClient, id string) ActionResult {
	var res ActionResult

	_, res.Err = client.Post(actionURL(client, id, "restart"),
        bytes.NewBufferString(""),
        &res.Body,
        nil)
	return res
}

// RebuildOptsBuilder is an interface that allows extensions to override the
// default behaviour of rebuild options
//type RebuildOptsBuilder interface {
//	ToServerRebuildMap() (map[string]interface{}, error)
//}

// RebuildOpts represents the configuration options used in a server rebuild
// operation
//type RebuildOpts struct {
//	// Required. The ID of the image you want your server to be provisioned on
//	ImageID string
//
//	// Name to set the server to
//	Name string
//
//	// Required. The server's admin password
//	AdminPass string
//
//	// AccessIPv4 [optional] provides a new IPv4 address for the instance.
//	AccessIPv4 string
//
//	// AccessIPv6 [optional] provides a new IPv6 address for the instance.
//	AccessIPv6 string
//
//	// Metadata [optional] contains key-value pairs (up to 255 bytes each) to attach to the server.
//	Metadata map[string]string
//
//	// Personality [optional] includes the path and contents of a file to inject into the server at launch.
//	// The maximum size of the file is 255 bytes (decoded).
//	Personality []byte
//}
//
//// ToServerRebuildMap formats a RebuildOpts struct into a map for use in JSON
//func (opts RebuildOpts) ToServerRebuildMap() (map[string]interface{}, error) {
//	var err error
//	server := make(map[string]interface{})
//
//	if opts.AdminPass == "" {
//		err = fmt.Errorf("AdminPass is required")
//	}
//
//	if opts.ImageID == "" {
//		err = fmt.Errorf("ImageID is required")
//	}
//
//	if err != nil {
//		return server, err
//	}
//
//	server["name"] = opts.Name
//	server["adminPass"] = opts.AdminPass
//	server["imageRef"] = opts.ImageID
//
//	if opts.AccessIPv4 != "" {
//		server["accessIPv4"] = opts.AccessIPv4
//	}
//
//	if opts.AccessIPv6 != "" {
//		server["accessIPv6"] = opts.AccessIPv6
//	}
//
//	if opts.Metadata != nil {
//		server["metadata"] = opts.Metadata
//	}
//
//	if opts.Personality != nil {
//		encoded := base64.StdEncoding.EncodeToString(opts.Personality)
//		server["personality"] = &encoded
//	}
//
//	return map[string]interface{}{"rebuild": server}, nil
//}

// Rebuild will reprovision the server according to the configuration options
// provided in the RebuildOpts struct.
//func Rebuild(client *speedycloud.ServiceClient, id string, opts RebuildOptsBuilder) RebuildResult {
//	var result RebuildResult
//
//	if id == "" {
//		result.Err = fmt.Errorf("ID is required")
//		return result
//	}
//
//	reqBody, err := opts.ToServerRebuildMap()
//	if err != nil {
//		result.Err = err
//		return result
//	}
//
//	_, result.Err = client.Post(actionURL(client, id), reqBody, &result.Body, nil)
//	return result
//}

// ResizeOptsBuilder is an interface that allows extensions to override the default structure of
// a Resize request.
//type ResizeOptsBuilder interface {
//	ToServerResizeMap() (map[string]interface{}, error)
//}
//
//// ResizeOpts represents the configuration options used to control a Resize operation.
//type ResizeOpts struct {
//	// FlavorRef is the ID of the flavor you wish your server to become.
//	FlavorRef string
//}
//
//// ToServerResizeMap formats a ResizeOpts as a map that can be used as a JSON request body for the
//// Resize request.
//func (opts ResizeOpts) ToServerResizeMap() (map[string]interface{}, error) {
//	resize := map[string]interface{}{
//		"flavorRef": opts.FlavorRef,
//	}
//
//	return map[string]interface{}{"resize": resize}, nil
//}

// Resize instructs the provider to change the flavor of the server.
// Note that this implies rebuilding it.
// Unfortunately, one cannot pass rebuild parameters to the resize function.
// When the resize completes, the server will be in RESIZE_VERIFY state.
// While in this state, you can explore the use of the new server's configuration.
// If you like it, call ConfirmResize() to commit the resize permanently.
// Otherwise, call RevertResize() to restore the old configuration.
//func Resize(client *speedycloud.ServiceClient, id string, opts ResizeOptsBuilder) ActionResult {
//	var res ActionResult
//	reqBody, err := opts.ToServerResizeMap()
//	if err != nil {
//		res.Err = err
//		return res
//	}
//
//	_, res.Err = client.Post(actionURL(client, id), reqBody, nil, nil)
//	return res
//}
//
//// ConfirmResize confirms a previous resize operation on a server.
//// See Resize() for more details.
//func ConfirmResize(client *speedycloud.ServiceClient, id string) ActionResult {
//	var res ActionResult
//
//	reqBody := map[string]interface{}{"confirmResize": nil}
//	_, res.Err = client.Post(actionURL(client, id), reqBody, nil, &speedycloud.RequestOpts{
//		OkCodes: []int{201, 202, 204},
//	})
//	return res
//}

// RevertResize cancels a previous resize operation on a server.
// See Resize() for more details.
//func RevertResize(client *speedycloud.ServiceClient, id string) ActionResult {
//	var res ActionResult
//	reqBody := map[string]interface{}{"revertResize": nil}
//	_, res.Err = client.Post(actionURL(client, id), reqBody, nil, nil)
//	return res
//}
//
//// RescueOptsBuilder is an interface that allows extensions to override the
//// default structure of a Rescue request.
//type RescueOptsBuilder interface {
//	ToServerRescueMap() (map[string]interface{}, error)
//}

// RescueOpts represents the configuration options used to control a Rescue
// option.
//type RescueOpts struct {
//	// AdminPass is the desired administrative password for the instance in
//	// RESCUE mode. If it's left blank, the server will generate a password.
//	AdminPass string
//}
//
//// ToServerRescueMap formats a RescueOpts as a map that can be used as a JSON
//// request body for the Rescue request.
//func (opts RescueOpts) ToServerRescueMap() (map[string]interface{}, error) {
//	server := make(map[string]interface{})
//	if opts.AdminPass != "" {
//		server["adminPass"] = opts.AdminPass
//	}
//	return map[string]interface{}{"rescue": server}, nil
//}
//
//// Rescue instructs the provider to place the server into RESCUE mode.
//func Rescue(client *speedycloud.ServiceClient, id string, opts RescueOptsBuilder) RescueResult {
//	var result RescueResult
//
//	if id == "" {
//		result.Err = fmt.Errorf("ID is required")
//		return result
//	}
//	reqBody, err := opts.ToServerRescueMap()
//	if err != nil {
//		result.Err = err
//		return result
//	}
//
//	_, result.Err = client.Post(actionURL(client, id), reqBody, &result.Body, &speedycloud.RequestOpts{
//		OkCodes: []int{200},
//	})
//
//	return result
//}
//
//// ResetMetadataOptsBuilder allows extensions to add additional parameters to the
//// Reset request.
//type ResetMetadataOptsBuilder interface {
//	ToMetadataResetMap() (map[string]interface{}, error)
//}
//
//// MetadataOpts is a map that contains key-value pairs.
//type MetadataOpts map[string]string
//
//// ToMetadataResetMap assembles a body for a Reset request based on the contents of a MetadataOpts.
//func (opts MetadataOpts) ToMetadataResetMap() (map[string]interface{}, error) {
//	return map[string]interface{}{"metadata": opts}, nil
//}
//
//// ToMetadataUpdateMap assembles a body for an Update request based on the contents of a MetadataOpts.
//func (opts MetadataOpts) ToMetadataUpdateMap() (map[string]interface{}, error) {
//	return map[string]interface{}{"metadata": opts}, nil
//}
//
//// ResetMetadata will create multiple new key-value pairs for the given server ID.
//// Note: Using this operation will erase any already-existing metadata and create
//// the new metadata provided. To keep any already-existing metadata, use the
//// UpdateMetadatas or UpdateMetadata function.
//func ResetMetadata(client *speedycloud.ServiceClient, id string, opts ResetMetadataOptsBuilder) ResetMetadataResult {
//	var res ResetMetadataResult
//	metadata, err := opts.ToMetadataResetMap()
//	if err != nil {
//		res.Err = err
//		return res
//	}
//	_, res.Err = client.Put(metadataURL(client, id), metadata, &res.Body, &speedycloud.RequestOpts{
//		OkCodes: []int{200},
//	})
//	return res
//}
//
//// Metadata requests all the metadata for the given server ID.
//func Metadata(client *speedycloud.ServiceClient, id string) GetMetadataResult {
//	var res GetMetadataResult
//	_, res.Err = client.Get(metadataURL(client, id), &res.Body, nil)
//	return res
//}
//
//// UpdateMetadataOptsBuilder allows extensions to add additional parameters to the
//// Create request.
//type UpdateMetadataOptsBuilder interface {
//	ToMetadataUpdateMap() (map[string]interface{}, error)
//}
//
//// UpdateMetadata updates (or creates) all the metadata specified by opts for the given server ID.
//// This operation does not affect already-existing metadata that is not specified
//// by opts.
//func UpdateMetadata(client *speedycloud.ServiceClient, id string, opts UpdateMetadataOptsBuilder) UpdateMetadataResult {
//	var res UpdateMetadataResult
//	metadata, err := opts.ToMetadataUpdateMap()
//	if err != nil {
//		res.Err = err
//		return res
//	}
//	_, res.Err = client.Post(metadataURL(client, id), metadata, &res.Body, &speedycloud.RequestOpts{
//		OkCodes: []int{200},
//	})
//	return res
//}
//
//// MetadatumOptsBuilder allows extensions to add additional parameters to the
//// Create request.
//type MetadatumOptsBuilder interface {
//	ToMetadatumCreateMap() (map[string]interface{}, string, error)
//}
//
//// MetadatumOpts is a map of length one that contains a key-value pair.
//type MetadatumOpts map[string]string
//
//// ToMetadatumCreateMap assembles a body for a Create request based on the contents of a MetadataumOpts.
//func (opts MetadatumOpts) ToMetadatumCreateMap() (map[string]interface{}, string, error) {
//	if len(opts) != 1 {
//		return nil, "", errors.New("CreateMetadatum operation must have 1 and only 1 key-value pair.")
//	}
//	metadatum := map[string]interface{}{"meta": opts}
//	var key string
//	for k := range metadatum["meta"].(MetadatumOpts) {
//		key = k
//	}
//	return metadatum, key, nil
//}
//
//// CreateMetadatum will create or update the key-value pair with the given key for the given server ID.
//func CreateMetadatum(client *speedycloud.ServiceClient, id string, opts MetadatumOptsBuilder) CreateMetadatumResult {
//	var res CreateMetadatumResult
//	metadatum, key, err := opts.ToMetadatumCreateMap()
//	if err != nil {
//		res.Err = err
//		return res
//	}
//
//	_, res.Err = client.Put(metadatumURL(client, id, key), metadatum, &res.Body, &speedycloud.RequestOpts{
//		OkCodes: []int{200},
//	})
//	return res
//}
//
//// Metadatum requests the key-value pair with the given key for the given server ID.
//func Metadatum(client *speedycloud.ServiceClient, id, key string) GetMetadatumResult {
//	var res GetMetadatumResult
//	_, res.Err = client.Request("GET", metadatumURL(client, id, key), speedycloud.RequestOpts{
//		JSONResponse: &res.Body,
//	})
//	return res
//}
//
//// DeleteMetadatum will delete the key-value pair with the given key for the given server ID.
//func DeleteMetadatum(client *speedycloud.ServiceClient, id, key string) DeleteMetadatumResult {
//	var res DeleteMetadatumResult
//	_, res.Err = client.Delete(metadatumURL(client, id, key), &speedycloud.RequestOpts{
//		JSONResponse: &res.Body,
//	})
//	return res
//}
//
//// ListAddresses makes a request against the API to list the servers IP addresses.
//func ListAddresses(client *speedycloud.ServiceClient, id string) pagination.Pager {
//	createPageFn := func(r pagination.PageResult) pagination.Page {
//		return AddressPage{pagination.SinglePageBase(r)}
//	}
//	return pagination.NewPager(client, listAddressesURL(client, id), createPageFn)
//}
//
//// ListAddressesByNetwork makes a request against the API to list the servers IP addresses
//// for the given network.
//func ListAddressesByNetwork(client *speedycloud.ServiceClient, id, network string) pagination.Pager {
//	createPageFn := func(r pagination.PageResult) pagination.Page {
//		return NetworkAddressPage{pagination.SinglePageBase(r)}
//	}
//	return pagination.NewPager(client, listAddressesByNetworkURL(client, id, network), createPageFn)
//}
