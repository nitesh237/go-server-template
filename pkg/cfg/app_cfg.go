package cfg

import (
	"time"

	enumspb "go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/worker"
)

// Application is a config Struct that holds common application config for cross cutting concerns
// e.g. logging, tracings, etc.
type Application struct {
	ServerPorts *ServerPorts
	Logging     *Logging
	Auth        *Auth
}

// config struct for TemporalWorkerApplication
type TemporalWorkerApplication struct {
	*Application
	Namespace                 string
	TaskQueue                 string
	WorkerOptions             *worker.Options
	WorkflowParamsList        WorkflowParamsList
	DefaultActivityParamsList ActivityParamsList
}

func (t *TemporalWorkerApplication) GetWorkerOptions() worker.Options {
	if t == nil || t.WorkerOptions == nil {
		return worker.Options{}
	}
	return *t.WorkerOptions
}

type ServerPorts struct {
	HttpPort int
}

// Logging holds all the parameters for tunning the logger
type Logging struct {
	EnableLoggingToFile bool
	LogPath             string
	MaxSizeInMBs        int // megabytes
	MaxBackups          int // There will be MaxBackups + 1 total files
}

type Auth struct {
	ConfigFilePath string
}

type HttpClient struct {

	// Transport layer configurations
	Transport struct {
		DialContext struct {
			// Timeout during connection setup
			Timeout time.Duration

			// KeepAlive specifies the interval between keep-alive
			// probes for an active network connection.
			// If zero, keep-alive probes are sent with a default value
			// (currently 15 seconds), if supported by the protocol and operating
			// system. Network protocols or operating systems that do
			// not support keep-alives ignore this field.
			// If negative, keep-alive probes are disabled.
			KeepAlive time.Duration
		}

		// TLSHandshakeTimeout specifies the maximum amount of time waiting to
		// wait for a TLS handshake. Zero means no timeout.
		TLSHandshakeTimeout time.Duration

		// MaxIdleConns controls the maximum number of idle (keep-alive)
		// connections across all hosts. Zero means no limit.
		MaxIdleConns int

		// IdleConnTimeout is the maximum amount of time an idle
		// (keep-alive) connection will remain idle before closing
		// itself.
		// Zero means no limit.
		IdleConnTimeout time.Duration

		// MaxIdleConnsPerHost, if non-zero, controls the maximum idle
		// (keep-alive) connections to keep per-host. If zero,
		// DefaultMaxIdleConnsPerHost is used.
		MaxIdleConnsPerHost int

		// MaxConnsPerHost optionally limits the total number of
		// connections per host, including connections in the dialing,
		// active, and idle states. On limit violation, dials will block.
		//
		// Zero means no limit.
		MaxConnsPerHost int

		// InsecureSkipVerify controls whether a client verifies the server's certificate chain and host name.
		InsecureSkipVerify bool
	}

	// Timeout specifies a time limit for requests made by this
	// Client. The timeout includes connection time, any
	// redirects, and reading the response body. The timer remains
	// running after Get, Head, Post, or Do return and will
	// interrupt reading of the Response.Body.
	//
	// A Timeout of zero means no timeout.
	Timeout time.Duration

	// RetryPolicy specifies the retry policy for the client.
	// Optional: default is no retry and is only used when initializing http.NewRetryableHttpClient
	RetryParams *RetryParams
}

type RetryParams struct {
	RegularInterval              *RegularInterval
	RegularIntervalWithJitter    *RegularIntervalWithJitter
	ExponentialBackOff           *ExponentialBackOff
	ExponentialBackOffWithJitter *ExponentialBackOffWithJitter
	RandomizedInterval           *RandomizedInterval
	Hybrid                       *Hybrid
}

type ActivityParams struct {
	// name of the activity for which config is defined
	ActivityName string

	// ScheduleToCloseTimeout - Total time that a workflow is willing to wait for Activity to complete.
	// ScheduleToCloseTimeout limits the total time of an Activity's execution including retries
	// 		(use StartToCloseTimeout to limit the time of a single attempt).
	// The zero value of this uses default value.
	// Either this option or StartToClose is required: Defaults to unlimited.
	ScheduleToCloseTimeout time.Duration
	// StartToCloseTimeout - Maximum time of a single Activity execution attempt.
	// Note that the Temporal Server doesn't detect Worker process failures directly. It relies on this timeout
	// to detect that an Activity that didn't complete on time. So this timeout should be as short as the longest
	// possible execution of the Activity body. Potentially long running Activities must specify HeartbeatTimeout
	// and call Activity.RecordHeartbeat(ctx, "my-heartbeat") periodically for timely failure detection.
	// If ScheduleToClose is not provided then this timeout is required: Defaults to the ScheduleToCloseTimeout value.
	StartToCloseTimeout time.Duration
	// HeartbeatTimeout - Heartbeat interval. Activity must call Activity.RecordHeartbeat(ctx, "my-heartbeat")
	// before this interval passes after the last heartbeat or the Activity starts.
	HeartbeatTimeout time.Duration
	RetryParams      *RetryParams
}

type ChildWorkflowParams struct {
	// name of the workflow for which config is defined
	WorkflowName string

	// WorkflowExecutionTimeout - The end to end timeout for the child workflow execution including retries
	// and continue as new.
	// Optional: defaults to unlimited.
	WorkflowExecutionTimeout time.Duration

	// WorkflowRunTimeout - The timeout for a single run of the child workflow execution. Each retry or
	// continue as new should obey this timeout. Use WorkflowExecutionTimeout to specify how long the parent
	// is willing to wait for the child completion.
	// Optional: defaults to WorkflowExecutionTimeout
	WorkflowRunTimeout time.Duration

	// WorkflowTaskTimeout - Maximum execution time of a single Workflow Task. In the majority of cases there is
	// no need to change this timeout. Note that this timeout is not related to the overall Workflow duration in
	// any way. It defines for how long the Workflow can get blocked in the case of a Workflow Worker crash.
	// Default is 10 seconds. Maximum value allowed by the Temporal Server is 1 minute.
	WorkflowTaskTimeout time.Duration

	// WaitForCancellation - Whether to wait for canceled child workflow to be ended (child workflow can be ended
	// as: completed/failed/timedout/terminated/canceled)
	// Optional: default false
	WaitForCancellation bool

	// RetryParams specify how to retry child workflow if error happens.
	// Optional: default is no retry
	RetryParams *RetryParams

	// ParentClosePolicy - Optional policy to decide what to do for the child.
	// Default is Terminate (if onboarded to this feature)
	ParentClosePolicy enumspb.ParentClosePolicy
}

type WorkflowParams struct {
	// name of the workflow for which config is defined
	WorkflowName string

	ActivityParamsList []*ActivityParams

	ChildWorkflowParamsList []*ChildWorkflowParams
}

func (p *WorkflowParams) GetActivityParamsMap() map[string]*ActivityParams {
	mp := map[string]*ActivityParams{}

	for _, ap := range p.ActivityParamsList {
		mp[ap.ActivityName] = ap
	}

	return mp
}

func (p *WorkflowParams) GetChildWorkflowParamsMap() map[string]*ChildWorkflowParams {
	mp := map[string]*ChildWorkflowParams{}

	for _, ap := range p.ChildWorkflowParamsList {
		mp[ap.WorkflowName] = ap
	}

	return mp
}

type WorkflowParamsList []*WorkflowParams

func (l WorkflowParamsList) GetWorkflowParamsMap() map[string]*WorkflowParams {
	mp := map[string]*WorkflowParams{}

	for _, wp := range l {
		mp[wp.WorkflowName] = wp
	}

	return mp
}

type ActivityParamsList []*ActivityParams

func (l ActivityParamsList) GetActivityParamsMap() map[string]*ActivityParams {
	mp := map[string]*ActivityParams{}

	for _, ap := range l {
		mp[ap.ActivityName] = ap
	}

	return mp
}
