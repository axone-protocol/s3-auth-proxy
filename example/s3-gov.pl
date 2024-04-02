:- discontiguous([title/2,partOf/2,chapter/1,section/1,subSection/1,article/1,paragraph/2]).

title('1', 'Service').
chapter('1').
  title('1.1', 'Service management'). partOf('1.1', '1').
  section('1.1').
    title('1.1.1', 'Use'). partOf('1.1.1', '1.1').
    subSection('1.1.1').
      title('1.1.1.1', 'Use'). partOf('1.1.1.1', '1.1.1').
      article('1.1.1.1').
        title('1.1.1.1.1', 'Service is allowed'). partOf('1.1.1.1.1', '1.1.1.1').
        paragraph('1.1.1.1.1', permitted) :- action('service:use').

        title('1.1.1.1.2', 'In a specific zone'). partOf('1.1.1.1.2', '1.1.1.1').
        paragraph('1.1.1.1.2', prohibited) :-
          zone(Z),
          Z \== 'did:key:zQ3shsLQeHXRgYHrRyTv9BhPLfgxHis8VowebBK2JMtMo8wzQ',
          action('service:use').
