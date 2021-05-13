package gov

var _ ProposalContent = (*TextProposal)(nil)

// Text Proposals for test
type TextProposal struct {
	Proposal
	Title       string `json:"title"`       //  Title of the proposal
	Description string `json:"description"` //  Description of the proposal
}

func NewTextProposal(title, description string) *TextProposal {
	tx := &TextProposal{
		Title:       title,
		Description: description,
	}
	return tx
}

// nolint
func (tp TextProposal) GetTitle() string           { return tp.Title }
func (tp TextProposal) GetDescription() string     { return tp.Description }
func (tp TextProposal) ProposalType() ProposalKind { return ProposalTypeText }
