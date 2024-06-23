# Steps to run

## Update config.yaml
create initial config.yaml file define the directories

## initial check
Run the initial check to find out details of the images (to run locally on the Linux server)
```
  mage dataCheck:t01_find_min_and_max_frame 0 # 0 is the index of the auc directory
```

A sample output is
```
Running target: DataCheck:T01_find_min_and_max_frame
frames total count:56746
first 5 frames:[15 30 45 60 75]
last 5 frames:[851130 851145 851160 851175 851190]
suggest report inerval:5674.50
```

Update the config.yaml file with the details

## Sample the images
Sample the frames from 10% to 100% every 10% to validate the OCR extraction
example:
```
mage dataCheck:t03_sample_frame 0
Running target: DataCheck:T03_sample_frame
------------frame:52155---------------
text extracted:| PY
ee pI
ES
LOT512 SALVADOR DAL! | Rhinocéros (recto); Etude pour Rhinocéros (verso) BID | $ 750,000
itemIndex:512, title:SALVADOR DAL! | Rhinocéros (recto); Etude pour Rhinocéros (verso) , price:750000.00, confidence:100
-----------------------
------------frame:104310---------------
text extracted:LOTS36 PAUL GAUGUIN | Etude pour La nativté tahitionne ‘ID | $50,000
itemIndex:, title:, price:0.00, confidence:0
-----------------------
------------frame:156465---------------
text extracted:crisis Ms Yo
itemIndex:, title:, price:0.00, confidence:0
-----------------------
------------frame:208620---------------
text extracted:Li
Pe | bull
if WY
3 | 六
rH ay
js
OTS89 EDGAR DEGAS | Etude de dansouse BD | $230,000
itemIndex:, title:, price:0.00, confidence:0
-----------------------
------------frame:260775---------------
text extracted:人 mame F |
itemIndex:, title:, price:0.00, confidence:0
-----------------------
------------frame:312915---------------
text extracted:
itemIndex:, title:, price:0.00, confidence:0
-----------------------
------------frame:365070---------------
text extracted:
itemIndex:, title:, price:0.00, confidence:0
-----------------------
------------frame:417225---------------
text extracted:
itemIndex:, title:, price:0.00, confidence:0
-----------------------
------------frame:469380---------------
text extracted:woh
at
BH
LOT 710. GUSTAVE LOISEAU | Bord dt vive Bio | $ 24,000
itemIndex:, title:, price:0.00, confidence:0
-----------------------
------------frame:521535---------------
text extracted:LOT 719. AUGUSTE RODIN | Ombre, taille originale dite taille de la Porto BID | $ 175,000
itemIndex:719, title:. AUGUSTE RODIN | Ombre, taille originale dite taille de la Porto , price:175000.00, confidence:100
-----------------------
```