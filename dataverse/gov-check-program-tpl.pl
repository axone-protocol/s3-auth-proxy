:- consult('{{ .GovCode }}').

action('{{ .Action }}').
subject('{{ .Subject }}').
zone('{{ .Zone }}').

tell(Result, Evidence) :-
    bagof(P:Modality, paragraph(P, Modality), Evidence),
    (   member(_: 'prohibited', Evidence) -> Result = 'prohibited'
    ;   member(_: 'permitted', Evidence) -> Result = 'permitted'
    ).
