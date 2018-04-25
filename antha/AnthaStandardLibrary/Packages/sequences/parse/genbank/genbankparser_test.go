package genbank

import (
	"fmt"
	"reflect"
	"testing"
)

type genbanktest struct {
	fileContents         []byte
	testname             string
	expectedfeaturenames []string
	featurePositionMap   map[string][2]int
}

var testfileContents string = `LOCUS       pRubiC-T2A-Cas9 11976 bp DNA SYN
DEFINITION  pRubiC-T2A-Cas9
ACCESSION   
KEYWORDS    
SOURCE      
ORGANISM  other sequences; artificial sequences; vectors.
FEATURES             Location/Qualifiers
     source          1..11976
                     /organism="pRubiC-T2A-Cas9"
                     /mol_type="other DNA"
     gene            complement(12..200)
                     /label="Ampicillin (860 - 672)"
                     /gene="Ampicillin (860 - 672)"
     promoter        complement(242..270)
                     /label="AmpR_promoter"
     promoter        492..1067
                     /label="CMV_immearly_promoter"
     misc_feature    556..1249
                     /label="5_LTR"
     misc_feature    571..858
                     /label="CAG_enhancer"
     misc_feature    1024..1044
                     /label="CMV_fwd_primer"
     misc_feature    1319..2127
                     /label="psi_plus_pack"
     misc_feature    1715..2121
                     /label="gag"
     CDS             1855..2334
                     /label="ORF frame 1"
                     /translation="MARDGTFNRGLITQVKIKVFSPGPHGHPDQVPYIVTWEALAFDP
                     PPWVKPFVHPKPPPPLPPSAPSLPLEPPRSTPPRSSLYPALTPSLGAGIEDLTQRPHT
                     LRTLINPCRLQIWPPRRVLAPPAGAPLLTASAATSDEGRSERPDPSARTLRTAARCS*"
     misc_feature    2057..2079
                     /label="MSCV_primer"
                     /translation="MARDGTFNRGLITQVKIKVFSPGPHGHPDQVPYIVTWEALAFDP
                     PPWVKPFVHPKPPPPLPPSAPSLPLEPPRSTPPRSSLYPALTPSLGAGIEDLTQRPHT
                     LRTLINPCRLQIWPPRRVLAPPAGAPLLTASAATSDEGRSERPDPSARTLRTAARCS*"
     promoter        2193..3408
                     /label="hUbC_promoter"
                     /translation="MARDGTFNRGLITQVKIKVFSPGPHGHPDQVPYIVTWEALAFDP
                     PPWVKPFVHPKPPPPLPPSAPSLPLEPPRSTPPRSSLYPALTPSLGAGIEDLTQRPHT
                     LRTLINPCRLQIWPPRRVLAPPAGAPLLTASAATSDEGRSERPDPSARTLRTAARCS*"
     ORF             3462..8498
                     /label="ORF frame 3"
                     /translation="MVSKGEEDNMAIIKEFMRFKVHMEGSVNGHEFEIEGEGEGRPYE
                     GTQTAKLKVTKGGPLPFAWDILSPQFMYGSKAYVKHPADIPDYLKLSFPEGFKWERVM
                     NFEDGGVVTVTQDSSLQDGEFIYKVKLRGTNFPSDGPVMQKKTMGWEASSERMYPEDG
                     ALKGEIKQRLKLKDGGHYDAEVKTTYKAKKPVQLPGAYNVNIKLDITSHNEDYTIVEQ
                     YERAEGRHSTGGMDELYKEGRGSLLTCGDVEENPGPEFDYKDHDGDYKDHDIDYKDDD
                     DKMAPKKKRKVGIHGVPAADKKYSIGLDIGTNSVGWAVITDEYKVPSKKFKVLGNTDR
                     HSIKKNLIGALLFDSGETAEATRLKRTARRRYTRRKNRICYLQEIFSNEMAKVDDSFF
                     HRLEESFLVEEDKKHERHPIFGNIVDEVAYHEKYPTIYHLRKKLVDSTDKADLRLIYL
                     ALAHMIKFRGHFLIEGDLNPDNSDVDKLFIQLVQTYNQLFEENPINASGVDAKAILSA
                     RLSKSRRLENLIAQLPGEKKNGLFGNLIALSLGLTPNFKSNFDLAEDAKLQLSKDTYD
                     DDLDNLLAQIGDQYADLFLAAKNLSDAILLSDILRVNTEITKAPLSASMIKRYDEHHQ
                     DLTLLKALVRQQLPEKYKEIFFDQSKNGYAGYIDGGASQEEFYKFIKPILEKMDGTEE
                     LLVKLNREDLLRKQRTFDNGSIPHQIHLGELHAILRRQEDFYPFLKDNREKIEKILTF
                     RIPYYVGPLARGNSRFAWMTRKSEETITPWNFEEVVDKGASAQSFIERMTNFDKNLPN
                     EKVLPKHSLLYEYFTVYNELTKVKYVTEGMRKPAFLSGEQKKAIVDLLFKTNRKVTVK
                     QLKEDYFKKIECFDSVEISGVEDRFNASLGTYHDLLKIIKDKDFLDNEENEDILEDIV
                     LTLTLFEDREMIEERLKTYAHLFDDKVMKQLKRRRYTGWGRLSRKLINGIRDKQSGKT
                     ILDFLKSDGFANRNFMQLIHDDSLTFKEDIQKAQVSGQGDSLHEHIANLAGSPAIKKG
                     ILQTVKVVDELVKVMGRHKPENIVIEMARENQTTQKGQKNSRERMKRIEEGIKELGSQ
                     ILKEHPVENTQLQNEKLYLYYLQNGRDMYVDQELDINRLSDYDVDHIVPQSFLKDDSI
                     DNKVLTRSDKNRGKSDNVPSEEVVKKMKNYWRQLLNAKLITQRKFDNLTKAERGGLSE
                     LDKAGFIKRQLVETRQITKHVAQILDSRMNTKYDENDKLIREVKVITLKSKLVSDFRK
                     DFQFYKVREINNYHHAHDAYLNAVVGTALIKKYPKLESEFVYGDYKVYDVRKMIAKSE
                     QEIGKATAKYFFYSNIMNFFKTEITLANGEIRKRPLIETNGETGEIVWDKGRDFATVR
                     KVLSMPQVNIVKKTEVQTGGFSKESILPKRNSDKLIARKKDWDPKKYGGFDSPTVAYS
                     VLVVAKVEKGKSKKLKSVKELLGITIMERSSFEKNPIDFLEAKGYKEVKKDLIIKLPK
                     YSLFELENGRKRMLASAGELQKGNELALPSKYVNFLYLASHYEKLKGSPEDNEQKQLF
                     VEQHKHYLDEIIEQISEFSKRVILADANLDKVLSAYNKHRDKPIREQAENIIHLFTLT
                     NLGAPAAFKYFDTTIDRKRYTSTKEVLDATLIHQSITGLYETRIDLSQLGGDKRPAAT
                     KKAGQAKKKK*"
     CDS             complement(4188..6863)
                     /label="ORF frame 2"
                     /translation="MVHIVVGQPVDVQFLVHVHIPPILQVVQVQLLVLQLGVFHGVFF
                     QDLAAQLFDALFDPLHSLAAVLLSLLGGLVLSGHFDHDVLGLVPAHHFHELVHHLHCL
                     QDALLNGGAAGQIGNVLVQAIALAGHLGFLDVLFKGQAVVVDQLHEVSVGEAVGLQEI
                     QDCLAGLLVPDAVDQLPAQPAPAGVSPPLQLLHHFVVEQVGIGFQPFLDHLSVLKQCQ
                     GQHDIFQNVLVFLIVQEVLVLDNFQQIVVCAQGGVEPIFHAGDFHGVEALDFLEVVLF
                     QLLHGHFPVGLEQQVHDGLFLLAAQEGGLSHSLGHVFHFGQLVIHGEVLVQQAVLGQH
                     LLVGQVLIEVGHPLDEALGGSALVHHFLEVPGGDGFLALSGHPGESAVSPGQRAHVVG
                     DAEGQDLLDLFPVVLQEWVKIFLPPQNGVQLSQVDLVGDAAVVEGPLLPQQVLSVQLH
                     EQFLGAVHLFQDGLDELVELFLAGSAVNVAGVAVLALVEENLFVLLRQLLPHESFQQG
                     QVLVVLVVSLDHRGAQGGLGDLGVHSQDVAQQDGVGQVLGGQKQVGVLVADLGQQVVQ
                     VVVVGVLAQLQFGILGQVEVALEVGGQAQAQGNQVSEQAILLLAGQLGDQIFQPSALA
                     QSGRQDGLGVHAAGVDGVFLEQLVVGLHQLDEQLVHVAVVGVQVALDQEVAPELDHVG
                     QGQIDQPQVGLVGAVHQFLSQVVDGGVLLVVGHLVHDVAEDGVPLVLLILFHQEGLFQ
                     SVEEAVVHLGHLVAEDLLQIADPVLPSGVSSSGGSLQPGGLGCFAAVEQQGSDQVLLD
                     AVPVGVAQHLEFLAGHLVLVGDHGPAHRVGADVQADAVLLVGCWDSVDTDLPLLLWGH
                     LIVIVFVINIMILVVSVVVLIVEFRAGILLHVTAC*"
ORIGIN
    1 atcattggaa aacgttcttc ggggcgaaaa ctctcaagga tcttaccgct gttgagatcc
   61 agttcgatgt aacccactcg tgcacccaac tgatcttcag catcttttac tttcaccagc
  121 gtttctgggt gagcaaaaac aggaaggcaa aatgccgcaa aaaagggaat aagggcgaca
  181 cggaaatgtt gaatactcat actcttcctt tttcaatatt attgaagcat ttatcagggt
  241 tattgtctca tgagcggata catatttgaa tgtatttaga aaaataaaca aataggggtt
  301 ccgcgcacat ttccccgaaa agtgccacct gacgtctaag aaaccattat tatcatgaca
  361 ttaacctata aaaataggcg tatcacgagg ccctttcgtc ttcaagaatt agcttggcca
  421 ttgcatacgt tgtatccata tcataatatg tacatttata ttggctcatg tccaacatta
  481 ccgccatgtt gacattgatt attgactagt tattaatagt aatcaattac ggggtcatta
  541 gttcatagcc catatatgga gttccgcgtt acataactta cggtaaatgg cccgcctggc
  601 tgaccgccca acgacccccg cccattgacg tcaataatga cgtatgttcc catagtaacg
  661 ccaataggga ctttccattg acgtcaatgg gtggagtatt tacggtaaac tgcccacttg
  721 gcagtacatc aagtgtatca tatgccaagt acgcccccta ttgacgtcaa tgacggtaaa
  781 tggcccgcct ggcattatgc ccagtacatg accttatggg actttcctac ttggcagtac
  841 atctacgtat tagtcatcgc tattaccatg gtgatgcggt tttggcagta catcaatggg
  901 cgtggatagc ggtttgactc acggggattt ccaagtctcc accccattga cgtcaatggg
  961 agtttgtttg gcaccaaaat caacgggact ttccaaaatg tcgtaacaac tccgccccat
 1021 tgacgcaaat gggcggtagg cgtgtacggt gggaggtcta tataagcaga gctcaataaa
 1081 agagcccaca acccctcact cggcgcgcca gtcttccgat agactgcgtc gcccgggtac
 1141 ccgtattccc aataaagcct cttgctgttt gcatccgaat cgtggtctcg ctgttccttg
 1201 ggagggtctc ctctgagtga ttgactaccc acgacggggg tctttcattt gggggctcgt
 1261 ccgggatttg gagacccctg cccagggacc accgacccac caccgggagg taagctggcc
 1321 agcaacttat ctgtgtctgt ccgattgtct agtgtctatg tttgatgtta tgcgcctgcg
 1381 tctgtactag ttagctaact agctctgtat ctggcggacc cgtggtggaa ctgacgagtt
 1441 ctgaacaccc ggccgcaacc ctgggagacg tcccagggac tttgggggcc gtttttgtgg
 1501 cccgacctga ggaagggagt cgatgtggaa tccgaccccg tcaggatatg tggttctggt
 1561 aggagacgag aacctaaaac agttcccgcc tccgtctgaa tttttgcttt cggtttggaa
 1621 ccgaagccgc gcgtcttgtc tgctgcagcg ctgcagcatc gttctgtgtt gtctctgtct
 1681 gactgtgttt ctgtatttgt ctgaaaatta gggccagact gttaccactc ccttaagttt
 1741 gaccttaggt cactggaaag atgtcgagcg gatcgctcac aaccagtcgg tagatgtcaa
 1801 gaagagacgt tgggttacct tctgctctgc agaatggcca acctttaacg tcggatggcg
 1861 cgagacggca cctttaaccg aggactcatc acccaggtta agatcaaggt cttttcacct
 1921 ggcccgcatg gacacccaga ccaggtcccg tacatcgtga cctgggaagc cttggctttt
 1981 gacccccctc cctgggtcaa gccctttgta caccctaagc ctccgcctcc tcttcctcca
 2041 tccgccccgt ctctccccct tgaacctcct cgttcgaccc cgcctcgatc ctccctttat
 2101 ccagccctca ctccttctct aggcgccgga attgaagatc tcacacagcg gccgcacaca
 2161 cttcgaacct taattaaccc gtgtcggctc cagatctggc ctccgcgccg ggttttggcg
 2221 cctcccgcgg gcgcccccct cctcacggcg agcgctgcca cgtcagacga agggcgcagc
 2281 gagcgtcctg atccttccgc ccggacgctc aggacagcgg cccgctgctc ataagactcg
 2341 gccttagaac cccagtatca gcagaaggac attttaggac gggacttggg tgactctagg
 2401 gcactggttt tctttccaga gagcggaaca ggcgaggaaa agtagtccct tctcggcgat
 2461 tctgcggagg gatctccgtg gggcggtgaa cgccgatgat tatataagga cgcgccgggt
 2521 gtggcacagc tagttccgtc gcagccggga tttgggtcgc ggttcttgtt tgtggatcgc
 2581 tgtgatcgtc acttggtgag tagcgggctg ctgggctggc cggggctttc gtggccgccg
 2641 ggccgctcgg tgggacggaa gcgtgtggag agaccgccaa gggctgtagt ctgggtccgc
 2701 gagcaaggtt gccctgaact gggggttggg gggagcgcac aaaatggcgg ctgttcccga
 2761 gtcttgaatg gaagacgctt gtaaggcggg ctgtgaggtc gttgaaacaa ggtggggggc
 2821 atggtgggcg gcaagaaccc aaggtcttga ggccttcgct aatgcgggaa agctcttatt
 2881 cgggtgagat gggctggggc accatctggg gaccctgacg tgaagtttgt cactgactgg
 2941 agaactcggg tttgtcgtct ggttgcgggg gcggcagtta tgcggtgccg ttgggcagtg
 3001 cacccgtacc tttgggagcg cgcgcctcgt cgtgtcgtga cgtcacccgt tctgttggct
 3061 tataatgcag ggtggggcca cctgccggta ggtgtgcggt aggcttttct ccgtcgcagg
 3121 acgcagggtt cgggcctagg gtaggctctc ctgaatcgac aggcgccgga cctctggtga
 3181 ggggagggat aagtgaggcg tcagtttctt tggtcggttt tatgtaccta tcttcttaag
 3241 tagctgaagc tccggttttg aactatgcgc tcggggttgg cgagtgtgtt ttgtgaagtt
 3301 ttttaggcac cttttgaaat gtaatcattt gggtcaatat gtaattttca gtgttagact
 3361 agtaaattgt ccgctaaatt ctggccgttt ttggcttttt tgttagacga agcttgggct
 3421 gcaggtcgac tctagaggat ccccgggtac cggtcgccac catggtgagc aagggcgagg
 3481 aggataacat ggccatcatc aaggagttca tgcgcttcaa ggtgcacatg gagggctccg
 3541 tgaacggcca cgagttcgag atcgagggcg agggcgaggg ccgcccctac gagggcaccc
 3601 agaccgccaa gctgaaggtg accaagggtg gccccctgcc cttcgcctgg gacatcctgt
 3661 cccctcagtt catgtacggc tccaaggcct acgtgaagca ccccgccgac atccccgact
 3721 acttgaagct gtccttcccc gagggcttca agtgggagcg cgtgatgaac ttcgaggacg
 3781 gcggcgtggt gaccgtgacc caggactcct ccctgcagga cggcgagttc atctacaagg
 3841 tgaagctgcg cggcaccaac ttcccctccg acggccccgt aatgcagaag aagaccatgg
 3901 gctgggaggc ctcctccgag cggatgtacc ccgaggacgg cgccctgaag ggcgagatca
 3961 agcagaggct gaagctgaag gacggcggcc actacgacgc tgaggtcaag accacctaca
 4021 aggccaagaa gcccgtgcag ctgcccggcg cctacaacgt caacatcaag ttggacatca
 4081 cctcccacaa cgaggactac accatcgtgg aacagtacga acgcgccgag ggccgccact
 4141 ccaccggcgg catggacgag ctgtacaagg agggcagagg aagtcttcta acatgcggtg
 4201 acgtggagga gaatcccggc cctgaattcg actataagga ccacgacgga gactacaagg
 4261 atcatgatat tgattacaaa gacgatgacg ataagatggc cccaaagaag aagcggaagg
 4321 tcggtatcca cggagtccca gcagccgaca agaagtacag catcggcctg gacatcggca
 4381 ccaactctgt gggctgggcc gtgatcaccg acgagtacaa ggtgcccagc aagaaattca
 4441 aggtgctggg caacaccgac cggcacagca tcaagaagaa cctgatcgga gccctgctgt
 4501 tcgacagcgg cgaaacagcc gaggccaccc ggctgaagag aaccgccaga agaagataca
 4561 ccagacggaa gaaccggatc tgctatctgc aagagatctt cagcaacgag atggccaagg
 4621 tggacgacag cttcttccac agactggaag agtccttcct ggtggaagag gataagaagc
 4681 acgagcggca ccccatcttc ggcaacatcg tggacgaggt ggcctaccac gagaagtacc
 4741 ccaccatcta ccacctgaga aagaaactgg tggacagcac cgacaaggcc gacctgcggc
 4801 tgatctatct ggccctggcc cacatgatca agttccgggg ccacttcctg atcgagggcg
 4861 acctgaaccc cgacaacagc gacgtggaca agctgttcat ccagctggtg cagacctaca
 4921 accagctgtt cgaggaaaac cccatcaacg ccagcggcgt ggacgccaag gccatcctgt
 4981 ctgccagact gagcaagagc agacggctgg aaaatctgat cgcccagctg cccggcgaga
 5041 agaagaatgg cctgttcgga aacctgattg ccctgagcct gggcctgacc cccaacttca
 5101 agagcaactt cgacctggcc gaggatgcca aactgcagct gagcaaggac acctacgacg
 5161 acgacctgga caacctgctg gcccagatcg gcgaccagta cgccgacctg tttctggccg
 5221 ccaagaacct gtccgacgcc atcctgctga gcgacatcct gagagtgaac accgagatca
 5281 ccaaggcccc cctgagcgcc tctatgatca agagatacga cgagcaccac caggacctga
 5341 ccctgctgaa agctctcgtg cggcagcagc tgcctgagaa gtacaaagag attttcttcg
 5401 accagagcaa gaacggctac gccggctaca ttgacggcgg agccagccag gaagagttct
 5461 acaagttcat caagcccatc ctggaaaaga tggacggcac cgaggaactg ctcgtgaagc
 5521 tgaacagaga ggacctgctg cggaagcagc ggaccttcga caacggcagc atcccccacc
 5581 agatccacct gggagagctg cacgccattc tgcggcggca ggaagatttt tacccattcc
 5641 tgaaggacaa ccgggaaaag atcgagaaga tcctgacctt ccgcatcccc tactacgtgg
 5701 gccctctggc caggggaaac agcagattcg cctggatgac cagaaagagc gaggaaacca
 5761 tcaccccctg gaacttcgag gaagtggtgg acaagggcgc ttccgcccag agcttcatcg
 5821 agcggatgac caacttcgat aagaacctgc ccaacgagaa ggtgctgccc aagcacagcc
 5881 tgctgtacga gtacttcacc gtgtataacg agctgaccaa agtgaaatac gtgaccgagg
 5941 gaatgagaaa gcccgccttc ctgagcggcg agcagaaaaa ggccatcgtg gacctgctgt
 6001 tcaagaccaa ccggaaagtg accgtgaagc agctgaaaga ggactacttc aagaaaatcg
 6061 agtgcttcga ctccgtggaa atctccggcg tggaagatcg gttcaacgcc tccctgggca
 6121 cataccacga tctgctgaaa attatcaagg acaaggactt cctggacaat gaggaaaacg
 6181 aggacattct ggaagatatc gtgctgaccc tgacactgtt tgaggacaga gagatgatcg
 6241 aggaacggct gaaaacctat gcccacctgt tcgacgacaa agtgatgaag cagctgaagc
 6301 ggcggagata caccggctgg ggcaggctga gccggaagct gatcaacggc atccgggaca
 6361 agcagtccgg caagacaatc ctggatttcc tgaagtccga cggcttcgcc aacagaaact
 6421 tcatgcagct gatccacgac gacagcctga cctttaaaga ggacatccag aaagcccagg
 6481 tgtccggcca gggcgatagc ctgcacgagc acattgccaa tctggccggc agccccgcca
 6541 ttaagaaggg catcctgcag acagtgaagg tggtggacga gctcgtgaaa gtgatgggcc
 6601 ggcacaagcc cgagaacatc gtgatcgaaa tggccagaga gaaccagacc acccagaagg
 6661 gacagaagaa cagccgcgag agaatgaagc ggatcgaaga gggcatcaaa gagctgggca
 6721 gccagatcct gaaagaacac cccgtggaaa acacccagct gcagaacgag aagctgtacc
 6781 tgtactacct gcagaatggg cgggatatgt acgtggacca ggaactggac atcaaccggc
 6841 tgtccgacta cgatgtggac catatcgtgc ctcagagctt tctgaaggac gactccatcg
 6901 acaacaaggt gctgaccaga agcgacaaga accggggcaa gagcgacaac gtgccctccg
 6961 aagaggtcgt gaagaagatg aagaactact ggcggcagct gctgaacgcc aagctgatta
 7021 cccagagaaa gttcgacaat ctgaccaagg ccgagagagg cggcctgagc gaactggata
 7081 aggccggctt catcaagaga cagctggtgg aaacccggca gatcacaaag cacgtggcac
 7141 agatcctgga ctcccggatg aacactaagt acgacgagaa tgacaagctg atccgggaag
 7201 tgaaagtgat caccctgaag tccaagctgg tgtccgattt ccggaaggat ttccagtttt
 7261 acaaagtgcg cgagatcaac aactaccacc acgcccacga cgcctacctg aacgccgtcg
 7321 tgggaaccgc cctgatcaaa aagtacccta agctggaaag cgagttcgtg tacggcgact
 7381 acaaggtgta cgacgtgcgg aagatgatcg ccaagagcga gcaggaaatc ggcaaggcta
 7441 ccgccaagta cttcttctac agcaacatca tgaacttttt caagaccgag attaccctgg
 7501 ccaacggcga gatccggaag cggcctctga tcgagacaaa cggcgaaacc ggggagatcg
 7561 tgtgggataa gggccgggat tttgccaccg tgcggaaagt gctgagcatg ccccaagtga
 7621 atatcgtgaa aaagaccgag gtgcagacag gcggcttcag caaagagtct atcctgccca
 7681 agaggaacag cgataagctg atcgccagaa agaaggactg ggaccctaag aagtacggcg
 7741 gcttcgacag ccccaccgtg gcctattctg tgctggtggt ggccaaagtg gaaaagggca
 7801 agtccaagaa actgaagagt gtgaaagagc tgctggggat caccatcatg gaaagaagca
 7861 gcttcgagaa gaatcccatc gactttctgg aagccaaggg ctacaaagaa gtgaaaaagg
 7921 acctgatcat caagctgcct aagtactccc tgttcgagct ggaaaacggc cggaagagaa
 7981 tgctggcctc tgccggcgaa ctgcagaagg gaaacgaact ggccctgccc tccaaatatg
 8041 tgaacttcct gtacctggcc agccactatg agaagctgaa gggctccccc gaggataatg
 8101 agcagaaaca gctgtttgtg gaacagcaca agcactacct ggacgagatc atcgagcaga
 8161 tcagcgagtt ctccaagaga gtgatcctgg ccgacgctaa tctggacaaa gtgctgtccg
 8221 cctacaacaa gcaccgggat aagcccatca gagagcaggc cgagaatatc atccacctgt
 8281 ttaccctgac caatctggga gcccctgccg ccttcaagta ctttgacacc accatcgacc
 8341 ggaagaggta caccagcacc aaagaggtgc tggacgccac cctgatccac cagagcatca
 8401 ccggcctgta cgagacacgg atcgacctgt ctcagctggg aggcgacaaa aggccggcgg
 8461 ccacgaaaaa ggccggccag gcaaaaaaga aaaagtaaga attcgatatc aagcttatcg
 8521 ataatcaacc tctggattac aaaatttgtg aaagattgac tggtattctt aactatgttg
 8581 ctccttttac gctatgtgga tacgctgctt taatgccttt gtatcatgct attgcttccc
 8641 gtatggcttt cattttctcc tccttgtata aatcctggtt gctgtctctt tatgaggagt
 8701 tgtggcccgt tgtcaggcaa cgtggcgtgg tgtgcactgt gtttgctgac gcaaccccca
 8761 ctggttgggg cattgccacc acctgtcagc tcctttccgg gactttcgct ttccccctcc
 8821 ctattgccac ggcggaactc atcgccgcct gccttgcccg ctgctggaca ggggctcggc
 8881 tgttgggcac tgacaattcc gtggtgttgt cggggaaatc atcgtccttt ccttggctgc
 8941 tcgcctgtgt tgccacctgg attctgcgcg ggacgtcctt ctgctacgtc ccttcggccc
 9001 tcaatccagc ggaccttcct tcccgcggcc tgctgccggc tctgcggcct cttccgcgtc
 9061 ttcgccttcg ccctcagacg agtcggatct ccctttgggc cgcctccccg catcgatacc
 9121 gtcgacctcg agatatcagt ggtccaggct ctagttttga ctcaacaata tcaccagctg
 9181 aagcctatag agtacgagcc atagataaaa taaaagattt tatttagtct ccagaaaaag
 9241 gggggaatga aagaccccac ctgtaggttt ggcaagctag cttaagtaac gccattttgc
 9301 aaggcatgga aaaatacata actgagaata gagaagttca gatcaaggtc aggaacagat
 9361 ggaacagggt cgaccctaga gaaccatcag atgtttccag ggtgccccaa ggacctgaaa
 9421 tgaccctgtg ccttatttga actaaccaat cagttcgctt ctcgcttctg ttcgcgcgct
 9481 tctgctcccc gagctcaata aaagagccca caacccctca ctcggggcgc cagtcctccg
 9541 attgactgag tcgcccgggt acccgtgtat ccaataaacc ctcttgcagt tgcatccgac
 9601 ttgtggtctc gctgttcctt gggagggtct cctctgagtg attgactacc cgtcagcggg
 9661 ggtctttcat ttgggggctc gtccgggatc gggagacccc tgcccaggga ccaccgaccc
 9721 accaccggga ggtaagctgg ctgcctcgcg cgtttcggtg atgacggtga aaacctctga
 9781 cacatgcagc tcccggagac ggtcacagct tgtctgtaag cggatgccgg gagcagacaa
 9841 gcccgtcagg gcgcgtcagc gggtgttggc gggtgtcggg gcgcagccat gacccagtca
 9901 cgtagcgata gcggagtgta gatccggctg tggaatgtgt gtcagttagg gtgtggaaag
 9961 tccccaggct ccccagcagg cagaagtatg caaagcatgc atctcaatta gtcagcaacc
10021 aggtgtgaaa gtccccaggc tcccagcagg cagaagtatg caaagcatgc atctcaatta
10081 gtcagcaacc atagtcccgc ccctactccg cccatcccgc ccctactccg cccagttccg
10141 cccaattctc cgcccccatg gctgactaat tttttttatt tatgcagagg ccgaggccgc
10201 ctcggcctct gagctattcc agaagtagtg aggaggcttt tttggaggcc taggcttttg
10261 caaaaagctt actggcttaa ctatgcggca tcagagcaga ttgtactgag agtgcaccat
10321 atgcggtgtg aaataccgca cagatgcgta aggagaaaat accgcatcag gcgctcttcc
10381 gcttcctcgc tcactgactc gctgcgctcg gtcgttcggc tgcggcgagc ggtatcagct
10441 cactcaaagg cggtaatacg gttatccaca gaatcagggg ataacgcagg aaagaacatg
10501 tgagcaaaag gccagcaaaa ggccaggaac cgtaaaaagg ccgcgttgct ggcgtttttc
10561 cataggctcc gcccccctga cgagcatcac aaaaatcgac gctcaagtca gaggtggcga
10621 aacccgacag gactataaag ataccaggcg tttccccctg gaagctccct cgtgcgctct
10681 cctgttccga ccctgccgct taccggatac ctgtccgcct ttctcccttc gggaagcgtg
10741 gcgctttctc atagctcacg ctgtaggtat ctcagttcgg tgtaggtcgt tcgctccaag
10801 ctgggctgtg tgcacgaacc ccccgttcag cccgaccgct gcgccttatc cggtaactat
10861 cgtcttgagt ccaacccggt aagacacgac ttatcgccac tggcagcagc cactggtaac
10921 aggattagca gagcgaggta tgtaggcggt gctacagagt tcttgaagtg gtggcctaac
10981 tacggctaca ctagaaggac agtatttggt atctgcgctc tgctgaagcc agttaccttc
11041 ggaaaaagag ttggtagctc ttgatccggc aaacaaacca ccgctggtag cggtggtttt
11101 tttgtttgca agcagcagat tacgcgcaga aaaaaaggat ctcaagaaga tcctttgatc
11161 ttttctacgg ggtctgacgc tcagtggaac gaaaactcac gttaagggat tttggtcatg
11221 agattatcaa aaaggatctt cacctagatc cttttaaatt aaaaatgaag ttttaaatca
11281 atctaaagta tatatgagta aacttggtct gacagttacc aatgcttaat cagtgaggca
11341 cctatctcag cgatctgtct atttcgttca tccatagttg cctgactccc cgtcgtgtag
11401 ataactacga tacgggaggg cttaccatct ggccccagtg ctgcaatgat accgcgagac
11461 ccacgctcac cggctccaga tttatcagca ataaaccagc cagccggaag ggccgagcgc
11521 agaagtggtc ctgcaacttt atccgcctcc atccagtcta ttaattgttg ccgggaagct
11581 agagtaagta gttcgccagt taatagtttg cgcaacgttg ttgccattgc tgcaggcatc
11641 gtggtgtcac gctcgtcgtt tggtatggct tcattcagct ccggttccca acgatcaagg
11701 cgagttacat gatcccccat gttgtgcaaa aaagcggtta gctccttcgg tcctccgatc
11761 gttgtcagaa gtaagttggc cgcagtgtta tcactcatgg ttatggcagc actgcataat
11821 tctcttactg tcatgccatc cgtaagatgc ttttctgtga ctggtgagta ctcaaccaag
11881 tcattctgag aatagtgtat gcggcgaccg agttgctctt gcccggcgtc aacacgggat
11941 aataccgcgc cacatagcag aactttaaaa gtgctc
//`

var tests []genbanktest = []genbanktest{
	{
		/*filename:             "test.gb",
		expectedfeaturenames: []string{"PluxI", "TetR", "LVA this is a test word", "RFP", "LVA", "TT BBa_B1002", "repA"},
		featurePositionMap: map[string][2]int{
			"PluxI": [2]int{582, 678},
			"TetR":  [2]int{752, 1417},
			"LVA this is a test word": [2]int{1373, 1411},
			"RFP":          [2]int{1484, 2203},
			"LVA":          [2]int{2159, 2197},
			"TT BBa_B1002": [2]int{2218, 2251},
			"repA":         [2]int{3971, 4921},
		},*/
		fileContents:         []byte(testfileContents),
		expectedfeaturenames: []string{"Ampicillin (860 - 672)", "Ampicillin (860 - 672)", "AmpR_promoter", "CMV_immearly_promoter", "5_LTR", "CAG_enhancer", "CMV_fwd_primer", "psi_plus_pack", "gag", "ORF frame 1", "MSCV_primer", "hUbC_promoter", "ORF frame 3", "ORF frame 2"},
		featurePositionMap: map[string][2]int{
			"Ampicillin (860 - 672)": {200, 12},
			"AmpR_promoter":          {270, 242},
			"CMV_immearly_promoter":  {492, 1067},
			"5_LTR":                  {556, 1249},
			"CAG_enhancer":           {571, 858},
			"CMV_fwd_primer":         {1024, 1044},
			"psi_plus_pack":          {1319, 2127},
			"gag":                    {1715, 2121},
			"ORF frame 1":            {1855, 2334},
			"MSCV_primer":            {2057, 2079},
			"hUbC_promoter":          {2193, 3408},
			"ORF frame 3":            {3462, 8498},
			"ORF frame 2":            {6863, 4188},
		},
	},
}

func TestGenbanktoAnnotatedSeq(t *testing.T) {

	for _, test := range tests {
		data := test.fileContents
		sequence, err := GenbankContentsToAnnotatedSeq(data)

		if err != nil {
			t.Error(
				"For", test.testname, "/n",
				"error: ", err.Error(), "\n",
			)
		}
		length := len(sequence.FeatureNames())
		length2 := len(test.expectedfeaturenames) //currently features that are identical are still appended this should be changed --> RemoveDuplicateFeatures in package search --> arrays
		if length != length2 {
			t.Error("ERROR: Not parsing correct number of elements.", length2, "features should be parsed but finding:", length)
		}
		if !reflect.DeepEqual(sequence.FeatureNames(), test.expectedfeaturenames) {
			t.Error(
				"NameERROR: For", test.testname, "\n",
				"expected", test.expectedfeaturenames, "\n",
				"got", sequence.FeatureNames(), "\n",
			)
		}
		for _, name := range sequence.FeatureNames() {
			features := sequence.GetFeatureByName(name)
			if len(features) > 0 {
				if len(features) > 1 {
					fmt.Println(
						"For", test.testname, "\n",
						"feature:", name, "\n",
						"expected positions:", test.featurePositionMap[name], "\n",
						"got more positions: ", len(features),
					)
				}
				for _, feature := range features {
					if feature.Start() != test.featurePositionMap[name][0] {
						t.Error(
							"NumberERROREnd: For", test.testname, "\n",
							"feature:", name, "\n",
							"expected", test.featurePositionMap[name][0], "\n",
							"got", feature.Start(), "\n IF 'expected 0' feature unspecified in featurePositionMap.",
						)
					}
					if feature.End() != test.featurePositionMap[name][1] {
						t.Error(
							"NumberERRORStart: For", test.testname, "\n",
							"feature:", name, "\n", "expected", test.featurePositionMap[name][1], "\n",
							"got", feature.End(), "\n IF 'expected 0' feature unspecified in featurePositionMap.",
						)
					}
				}
			}

		}
	}

}
