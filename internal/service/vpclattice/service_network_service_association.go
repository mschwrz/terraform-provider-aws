package vpclattice

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/vpclattice"
	"github.com/aws/aws-sdk-go-v2/service/vpclattice/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_vpclattice_service_network_service_association")
func ResourceServiceNetworkServiceAssociation() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceServiceNetworkServiceAssociationCreate,
		ReadWithoutTimeout:   resourceServiceNetworkServiceAssociationRead,
		UpdateWithoutTimeout: resourceServiceNetworkServiceAssociationUpdate,
		DeleteWithoutTimeout: resourceServiceNetworkServiceAssociationDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dns_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hosted_zone_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_identifier": {
				Type:     schema.TypeString,
				Required: true,
			},
			"service_network_identifier": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			names.AttrTags: tftags.TagsSchema(),
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

const (
	ResNameServiceNetworkServiceAssociation = "Service Network Service Association"
)

func resourceServiceNetworkServiceAssociationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// TIP: ==== RESOURCE CREATE ====
	// Generally, the Create function should do the following things. Make
	// sure there is a good reason if you don't do one of these.
	//
	// 1. Get a client connection to the relevant service
	// 2. Populate a create input structure
	// 3. Call the AWS create/put function
	// 4. Using the output from the create function, set the minimum arguments
	//    and attributes for the Read function to work. At a minimum, set the
	//    resource ID. E.g., d.SetId(<Identifier, such as AWS ID or ARN>)
	// 5. Use a waiter to wait for create to complete
	// 6. Call the Read function in the Create return

	// TIP: -- 1. Get a client connection to the relevant service
	conn := meta.(*conns.AWSClient).VPCLatticeClient()

	// TIP: -- 2. Populate a create input structure
	in := &vpclattice.CreateServiceNetworkServiceAssociationInput{
		// TIP: Mandatory or fields that will always be present can be set when
		// you create the Input structure. (Replace these with real fields.)
		ServiceNetworkServiceAssociationName: aws.String(d.Get("name").(string)),
		ServiceNetworkServiceAssociationType: aws.String(d.Get("type").(string)),
	}

	if v, ok := d.GetOk("max_size"); ok {
		// TIP: Optional fields should be set based on whether or not they are
		// used.
		in.MaxSize = aws.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("complex_argument"); ok && len(v.([]interface{})) > 0 {
		// TIP: Use an expander to assign a complex argument.
		in.ComplexArguments = expandComplexArguments(v.([]interface{}))
	}

	// TIP: Not all resources support tags and tags don't always make sense. If
	// your resource doesn't need tags, you can remove the tags lines here and
	// below. Many resources do include tags so this a reminder to include them
	// where possible.
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	tags := defaultTagsConfig.MergeTags(tftags.New(ctx, d.Get("tags").(map[string]interface{})))

	if len(tags) > 0 {
		in.Tags = Tags(tags.IgnoreAWS())
	}

	// TIP: -- 3. Call the AWS create function
	out, err := conn.CreateServiceNetworkServiceAssociation(ctx, in)
	if err != nil {
		// TIP: Since d.SetId() has not been called yet, you cannot use d.Id()
		// in error messages at this point.
		return create.DiagError(names.VPCLattice, create.ErrActionCreating, ResNameServiceNetworkServiceAssociation, d.Get("name").(string), err)
	}

	if out == nil || out.ServiceNetworkServiceAssociation == nil {
		return create.DiagError(names.VPCLattice, create.ErrActionCreating, ResNameServiceNetworkServiceAssociation, d.Get("name").(string), errors.New("empty output"))
	}

	// TIP: -- 4. Set the minimum arguments and/or attributes for the Read function to
	// work.
	d.SetId(aws.ToString(out.ServiceNetworkServiceAssociation.ServiceNetworkServiceAssociationID))

	// TIP: -- 5. Use a waiter to wait for create to complete
	if _, err := waitServiceNetworkServiceAssociationCreated(ctx, conn, d.Id(), d.Timeout(schema.TimeoutCreate)); err != nil {
		return create.DiagError(names.VPCLattice, create.ErrActionWaitingForCreation, ResNameServiceNetworkServiceAssociation, d.Id(), err)
	}

	// TIP: -- 6. Call the Read function in the Create return
	return resourceServiceNetworkServiceAssociationRead(ctx, d, meta)
}

func resourceServiceNetworkServiceAssociationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// TIP: ==== RESOURCE READ ====
	// Generally, the Read function should do the following things. Make
	// sure there is a good reason if you don't do one of these.
	//
	// 1. Get a client connection to the relevant service
	// 2. Get the resource from AWS
	// 3. Set ID to empty where resource is not new and not found
	// 4. Set the arguments and attributes
	// 5. Set the tags
	// 6. Return nil

	// TIP: -- 1. Get a client connection to the relevant service
	conn := meta.(*conns.AWSClient).VPCLatticeClient()

	// TIP: -- 2. Get the resource from AWS using an API Get, List, or Describe-
	// type function, or, better yet, using a finder.
	out, err := findServiceNetworkServiceAssociationByID(ctx, conn, d.Id())

	// TIP: -- 3. Set ID to empty where resource is not new and not found
	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] VPCLattice ServiceNetworkServiceAssociation (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return create.DiagError(names.VPCLattice, create.ErrActionReading, ResNameServiceNetworkServiceAssociation, d.Id(), err)
	}

	// TIP: -- 4. Set the arguments and attributes
	//
	// For simple data types (i.e., schema.TypeString, schema.TypeBool,
	// schema.TypeInt, and schema.TypeFloat), a simple Set call (e.g.,
	// d.Set("arn", out.Arn) is sufficient. No error or nil checking is
	// necessary.
	//
	// However, there are some situations where more handling is needed.
	// a. Complex data types (e.g., schema.TypeList, schema.TypeSet)
	// b. Where errorneous diffs occur. For example, a schema.TypeString may be
	//    a JSON. AWS may return the JSON in a slightly different order but it
	//    is equivalent to what is already set. In that case, you may check if
	//    it is equivalent before setting the different JSON.
	d.Set("arn", out.Arn)
	d.Set("name", out.Name)

	// TIP: Setting a complex type.
	// For more information, see:
	// https://hashicorp.github.io/terraform-provider-aws/data-handling-and-conversion/#data-handling-and-conversion
	// https://hashicorp.github.io/terraform-provider-aws/data-handling-and-conversion/#flatten-functions-for-blocks
	// https://hashicorp.github.io/terraform-provider-aws/data-handling-and-conversion/#root-typeset-of-resource-and-aws-list-of-structure
	if err := d.Set("complex_argument", flattenComplexArguments(out.ComplexArguments)); err != nil {
		return create.DiagError(names.VPCLattice, create.ErrActionSetting, ResNameServiceNetworkServiceAssociation, d.Id(), err)
	}

	// TIP: Setting a JSON string to avoid errorneous diffs.
	p, err := verify.SecondJSONUnlessEquivalent(d.Get("policy").(string), aws.ToString(out.Policy))
	if err != nil {
		return create.DiagError(names.VPCLattice, create.ErrActionSetting, ResNameServiceNetworkServiceAssociation, d.Id(), err)
	}

	p, err = structure.NormalizeJsonString(p)
	if err != nil {
		return create.DiagError(names.VPCLattice, create.ErrActionSetting, ResNameServiceNetworkServiceAssociation, d.Id(), err)
	}

	d.Set("policy", p)

	// TIP: -- 5. Set the tags
	//
	// TIP: Not all resources support tags and tags don't always make sense. If
	// your resource doesn't need tags, you can remove the tags lines here and
	// below. Many resources do include tags so this a reminder to include them
	// where possible.
	tags, err := ListTags(ctx, conn, d.Id())
	if err != nil {
		return create.DiagError(names.VPCLattice, create.ErrActionReading, ResNameServiceNetworkServiceAssociation, d.Id(), err)
	}

	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig
	tags = tags.IgnoreAWS().IgnoreConfig(ignoreTagsConfig)

	if err := d.Set("tags", tags.RemoveDefaultConfig(defaultTagsConfig).Map()); err != nil {
		return create.DiagError(names.VPCLattice, create.ErrActionSetting, ResNameServiceNetworkServiceAssociation, d.Id(), err)
	}

	if err := d.Set("tags_all", tags.Map()); err != nil {
		return create.DiagError(names.VPCLattice, create.ErrActionSetting, ResNameServiceNetworkServiceAssociation, d.Id(), err)
	}

	// TIP: -- 6. Return nil
	return nil
}

func resourceServiceNetworkServiceAssociationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// TIP: ==== RESOURCE UPDATE ====
	// Not all resources have Update functions. There are a few reasons:
	// a. The AWS API does not support changing a resource
	// b. All arguments have ForceNew: true, set
	// c. The AWS API uses a create call to modify an existing resource
	//
	// In the cases of a. and b., the main resource function will not have a
	// UpdateWithoutTimeout defined. In the case of c., Update and Create are
	// the same.
	//
	// The rest of the time, there should be an Update function and it should
	// do the following things. Make sure there is a good reason if you don't
	// do one of these.
	//
	// 1. Get a client connection to the relevant service
	// 2. Populate a modify input structure and check for changes
	// 3. Call the AWS modify/update function
	// 4. Use a waiter to wait for update to complete
	// 5. Call the Read function in the Update return

	// TIP: -- 1. Get a client connection to the relevant service
	conn := meta.(*conns.AWSClient).VPCLatticeClient()

	// TIP: -- 2. Populate a modify input structure and check for changes
	//
	// When creating the input structure, only include mandatory fields. Other
	// fields are set as needed. You can use a flag, such as update below, to
	// determine if a certain portion of arguments have been changed and
	// whether to call the AWS update function.
	update := false

	in := &vpclattice.UpdateServiceNetworkServiceAssociationInput{
		Id: aws.String(d.Id()),
	}

	if d.HasChanges("an_argument") {
		in.AnArgument = aws.String(d.Get("an_argument").(string))
		update = true
	}

	if !update {
		// TIP: If update doesn't do anything at all, which is rare, you can
		// return nil. Otherwise, return a read call, as below.
		return nil
	}

	// TIP: -- 3. Call the AWS modify/update function
	log.Printf("[DEBUG] Updating VPCLattice ServiceNetworkServiceAssociation (%s): %#v", d.Id(), in)
	out, err := conn.UpdateServiceNetworkServiceAssociation(ctx, in)
	if err != nil {
		return create.DiagError(names.VPCLattice, create.ErrActionUpdating, ResNameServiceNetworkServiceAssociation, d.Id(), err)
	}

	// TIP: -- 4. Use a waiter to wait for update to complete
	if _, err := waitServiceNetworkServiceAssociationUpdated(ctx, conn, aws.ToString(out.OperationId), d.Timeout(schema.TimeoutUpdate)); err != nil {
		return create.DiagError(names.VPCLattice, create.ErrActionWaitingForUpdate, ResNameServiceNetworkServiceAssociation, d.Id(), err)
	}

	// TIP: -- 5. Call the Read function in the Update return
	return resourceServiceNetworkServiceAssociationRead(ctx, d, meta)
}

func resourceServiceNetworkServiceAssociationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// TIP: ==== RESOURCE DELETE ====
	// Most resources have Delete functions. There are rare situations
	// where you might not need a delete:
	// a. The AWS API does not provide a way to delete the resource
	// b. The point of your resource is to perform an action (e.g., reboot a
	//    server) and deleting serves no purpose.
	//
	// The Delete function should do the following things. Make sure there
	// is a good reason if you don't do one of these.
	//
	// 1. Get a client connection to the relevant service
	// 2. Populate a delete input structure
	// 3. Call the AWS delete function
	// 4. Use a waiter to wait for delete to complete
	// 5. Return nil

	// TIP: -- 1. Get a client connection to the relevant service
	conn := meta.(*conns.AWSClient).VPCLatticeClient()

	// TIP: -- 2. Populate a delete input structure
	log.Printf("[INFO] Deleting VPCLattice ServiceNetworkServiceAssociation %s", d.Id())

	// TIP: -- 3. Call the AWS delete function
	_, err := conn.DeleteServiceNetworkServiceAssociation(ctx, &vpclattice.DeleteServiceNetworkServiceAssociationInput{
		Id: aws.String(d.Id()),
	})

	// TIP: On rare occassions, the API returns a not found error after deleting a
	// resource. If that happens, we don't want it to show up as an error.
	if err != nil {
		var nfe *types.ResourceNotFoundException
		if errors.As(err, &nfe) {
			return nil
		}

		return create.DiagError(names.VPCLattice, create.ErrActionDeleting, ResNameServiceNetworkServiceAssociation, d.Id(), err)
	}

	// TIP: -- 4. Use a waiter to wait for delete to complete
	if _, err := waitServiceNetworkServiceAssociationDeleted(ctx, conn, d.Id(), d.Timeout(schema.TimeoutDelete)); err != nil {
		return create.DiagError(names.VPCLattice, create.ErrActionWaitingForDeletion, ResNameServiceNetworkServiceAssociation, d.Id(), err)
	}

	// TIP: -- 5. Return nil
	return nil
}

// TIP: ==== STATUS CONSTANTS ====
// Create constants for states and statuses if the service does not
// already have suitable constants. We prefer that you use the constants
// provided in the service if available (e.g., amp.WorkspaceStatusCodeActive).
const (
	statusChangePending = "Pending"
	statusDeleting      = "Deleting"
	statusNormal        = "Normal"
	statusUpdated       = "Updated"
)

// TIP: ==== WAITERS ====
// Some resources of some services have waiters provided by the AWS API.
// Unless they do not work properly, use them rather than defining new ones
// here.
//
// Sometimes we define the wait, status, and find functions in separate
// files, wait.go, status.go, and find.go. Follow the pattern set out in the
// service and define these where it makes the most sense.
//
// If these functions are used in the _test.go file, they will need to be
// exported (i.e., capitalized).
//
// You will need to adjust the parameters and names to fit the service.

func waitServiceNetworkServiceAssociationCreated(ctx context.Context, conn *vpclattice.Client, id string, timeout time.Duration) (*vpclattice.ServiceNetworkServiceAssociation, error) {
	stateConf := &retry.StateChangeConf{
		Pending:                   []string{},
		Target:                    []string{statusNormal},
		Refresh:                   statusServiceNetworkServiceAssociation(ctx, conn, id),
		Timeout:                   timeout,
		NotFoundChecks:            20,
		ContinuousTargetOccurence: 2,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)
	if out, ok := outputRaw.(*vpclattice.ServiceNetworkServiceAssociation); ok {
		return out, err
	}

	return nil, err
}

// TIP: It is easier to determine whether a resource is updated for some
// resources than others. The best case is a status flag that tells you when
// the update has been fully realized. Other times, you can check to see if a
// key resource argument is updated to a new value or not.

func waitServiceNetworkServiceAssociationUpdated(ctx context.Context, conn *vpclattice.Client, id string, timeout time.Duration) (*vpclattice.ServiceNetworkServiceAssociation, error) {
	stateConf := &retry.StateChangeConf{
		Pending:                   []string{statusChangePending},
		Target:                    []string{statusUpdated},
		Refresh:                   statusServiceNetworkServiceAssociation(ctx, conn, id),
		Timeout:                   timeout,
		NotFoundChecks:            20,
		ContinuousTargetOccurence: 2,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)
	if out, ok := outputRaw.(*vpclattice.ServiceNetworkServiceAssociation); ok {
		return out, err
	}

	return nil, err
}

// TIP: A deleted waiter is almost like a backwards created waiter. There may
// be additional pending states, however.

func waitServiceNetworkServiceAssociationDeleted(ctx context.Context, conn *vpclattice.Client, id string, timeout time.Duration) (*vpclattice.ServiceNetworkServiceAssociation, error) {
	stateConf := &retry.StateChangeConf{
		Pending: []string{statusDeleting, statusNormal},
		Target:  []string{},
		Refresh: statusServiceNetworkServiceAssociation(ctx, conn, id),
		Timeout: timeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)
	if out, ok := outputRaw.(*vpclattice.ServiceNetworkServiceAssociation); ok {
		return out, err
	}

	return nil, err
}

// TIP: ==== STATUS ====
// The status function can return an actual status when that field is
// available from the API (e.g., out.Status). Otherwise, you can use custom
// statuses to communicate the states of the resource.
//
// Waiters consume the values returned by status functions. Design status so
// that it can be reused by a create, update, and delete waiter, if possible.

func statusServiceNetworkServiceAssociation(ctx context.Context, conn *vpclattice.Client, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		out, err := findServiceNetworkServiceAssociationByID(ctx, conn, id)
		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return out, aws.ToString(out.Status), nil
	}
}

// TIP: ==== FINDERS ====
// The find function is not strictly necessary. You could do the API
// request from the status function. However, we have found that find often
// comes in handy in other places besides the status function. As a result, it
// is good practice to define it separately.

func findServiceNetworkServiceAssociationByID(ctx context.Context, conn *vpclattice.Client, id string) (*vpclattice.ServiceNetworkServiceAssociation, error) {
	in := &vpclattice.GetServiceNetworkServiceAssociationInput{
		Id: aws.String(id),
	}
	out, err := conn.GetServiceNetworkServiceAssociation(ctx, in)
	if err != nil {
		var nfe *types.ResourceNotFoundException
		if errors.As(err, &nfe) {
			return nil, &retry.NotFoundError{
				LastError:   err,
				LastRequest: in,
			}
		}

		return nil, err
	}

	if out == nil || out.ServiceNetworkServiceAssociation == nil {
		return nil, tfresource.NewEmptyResultError(in)
	}

	return out.ServiceNetworkServiceAssociation, nil
}

// TIP: ==== FLEX ====
// Flatteners and expanders ("flex" functions) help handle complex data
// types. Flatteners take an API data type and return something you can use in
// a d.Set() call. In other words, flatteners translate from AWS -> Terraform.
//
// On the other hand, expanders take a Terraform data structure and return
// something that you can send to the AWS API. In other words, expanders
// translate from Terraform -> AWS.
//
// See more:
// https://hashicorp.github.io/terraform-provider-aws/data-handling-and-conversion/
func flattenComplexArgument(apiObject *vpclattice.ComplexArgument) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	m := map[string]interface{}{}

	if v := apiObject.SubFieldOne; v != nil {
		m["sub_field_one"] = aws.ToString(v)
	}

	if v := apiObject.SubFieldTwo; v != nil {
		m["sub_field_two"] = aws.ToString(v)
	}

	return m
}

// TIP: Often the AWS API will return a slice of structures in response to a
// request for information. Sometimes you will have set criteria (e.g., the ID)
// that means you'll get back a one-length slice. This plural function works
// brilliantly for that situation too.
func flattenComplexArguments(apiObjects []*vpclattice.ComplexArgument) []interface{} {
	if len(apiObjects) == 0 {
		return nil
	}

	var l []interface{}

	for _, apiObject := range apiObjects {
		if apiObject == nil {
			continue
		}

		l = append(l, flattenComplexArgument(apiObject))
	}

	return l
}

// TIP: Remember, as mentioned above, expanders take a Terraform data structure
// and return something that you can send to the AWS API. In other words,
// expanders translate from Terraform -> AWS.
//
// See more:
// https://hashicorp.github.io/terraform-provider-aws/data-handling-and-conversion/
func expandComplexArgument(tfMap map[string]interface{}) *vpclattice.ComplexArgument {
	if tfMap == nil {
		return nil
	}

	a := &vpclattice.ComplexArgument{}

	if v, ok := tfMap["sub_field_one"].(string); ok && v != "" {
		a.SubFieldOne = aws.String(v)
	}

	if v, ok := tfMap["sub_field_two"].(string); ok && v != "" {
		a.SubFieldTwo = aws.String(v)
	}

	return a
}

// TIP: Even when you have a list with max length of 1, this plural function
// works brilliantly. However, if the AWS API takes a structure rather than a
// slice of structures, you will not need it.
func expandComplexArguments(tfList []interface{}) []*vpclattice.ComplexArgument {
	// TIP: The AWS API can be picky about whether you send a nil or zero-
	// length for an argument that should be cleared. For example, in some
	// cases, if you send a nil value, the AWS API interprets that as "make no
	// changes" when what you want to say is "remove everything." Sometimes
	// using a zero-length list will cause an error.
	//
	// As a result, here are two options. Usually, option 1, nil, will work as
	// expected, clearing the field. But, test going from something to nothing
	// to make sure it works. If not, try the second option.

	// TIP: Option 1: Returning nil for zero-length list
	if len(tfList) == 0 {
		return nil
	}

	var s []*vpclattice.ComplexArgument

	// TIP: Option 2: Return zero-length list for zero-length list. If option 1 does
	// not work, after testing going from something to nothing (if that is
	// possible), uncomment out the next line and remove option 1.
	//
	// s := make([]*vpclattice.ComplexArgument, 0)

	for _, r := range tfList {
		m, ok := r.(map[string]interface{})

		if !ok {
			continue
		}

		a := expandComplexArgument(m)

		if a == nil {
			continue
		}

		s = append(s, a)
	}

	return s
}
