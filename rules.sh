#!sh

F1=${F1:-NOUN}
F2=${F2:-VERB}
RATIO=${RATIO:-10}
MIN=${MIN:-500}

sqlite3 rules.db "
SELECT length(feature),feature,key, '$F1' AS resolve, round((1.0*$F1)/$F2,1) AS ratio FROM features
  WHERE $F1+$F2>$MIN AND (1.0*$F1)/$F2>$RATIO
  UNION ALL
  SELECT length(feature),feature, key, '$F2' AS resolve, round((1.0*$F2)/$F1,1) as ratio FROM features
  WHERE $F1+$F2>$MIN AND (1.0*$F2)/$F1>$RATIO
  ORDER BY length(feature) DESC, ratio DESC;
"
