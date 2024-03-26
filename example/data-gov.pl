:- discontiguous([title/2,partOf/2,chapter/1,section/1,subSection/1,article/1,paragraph/2]).

title('1', 'Data').
chapter('1').
  title('1.1', 'Dataset management'). partOf('1.1', '1').
  section('1.1').
    title('1.1.1', 'Access'). partOf('1.1.1', '1.1').
    subSection('1.1.1').
      title('1.1.1.1', 'Read'). partOf('1.1.1.1', '1.1.1').
      article('1.1.1.1').
        title('1.1.1.1.1', 'Read access is allowed'). partOf('1.1.1.1.1', '1.1.1.1').
        paragraph('1.1.1.1.1', permitted) :- action('dataset:read').

        title('1.1.1.1.2', 'For only permitted identities'). partOf('1.1.1.1.2', '1.1.1.1').
        paragraph('1.1.1.1.2', prohibited) :-
          subject(S),
    	  S \== 'did:key:zQ3shuJ8dEheKxgedw1ZQvP1PHthG2vokpCTdqggHvBxZ55Bz',
          action('dataset:read').

        title('1.1.1.1.3', 'In a specific zone'). partOf('1.1.1.1.3', '1.1.1.1').
        paragraph('1.1.1.1.3', prohibited) :-
          zone(Z),
          Z \== 'did:key:zQ3shsLQeHXRgYHrRyTv9BhPLfgxHis8VowebBK2JMtMo8wzQ',
          action('dataset:read').
