package workflow

type State string

const (
	DraftState    State = "DRAFT"
	PreviewState  State = "PREVIEW"
	OnlineState   State = "ONLINE"
	OfflineState  State = "OFFLINE"
	ArchivedState State = "ARCHIVED"
)

func (state State) ToWorkflow() Workflow {
	switch state {
	case DraftState:
		return &Draft{}
	case PreviewState:
		return &Preview{}
	case OnlineState:
		return &Online{}
	case OfflineState:
		return &Offline{}
	case ArchivedState:
		return &Archived{}
	default:
		return nil
	}
}
