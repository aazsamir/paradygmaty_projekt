:- set_prolog_flag(verbose, silent).
:- initialization main.
head([H|_], H).
main :-
    current_prolog_flag(argv, Argv),
    head(Argv, Who),
    atom_number(Who, WhoNum),
    related(WhoNum, X),
    format('~w', X),
    halt.
main :-
    format('0'),
    halt.

related(1, 1).

related(2, 1).
