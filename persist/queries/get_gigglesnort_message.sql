-- ?1: the word that is maybe a gigglesnort word
-- returns zero or one row (text): the gigglesnort message

SELECT message FROM gigglesnort WHERE word = ?1;
