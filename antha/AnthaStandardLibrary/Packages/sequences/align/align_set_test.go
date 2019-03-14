package align

import (
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences/parse/fasta"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"io/ioutil"
	"testing"
)

func TestDNASet(t *testing.T) {

	fastaText := `
>CCE57619 cdna plasmid:HUSEC2011CHR1:pHUSEC2011-2:1219:1338:-1 gene:HUS2011_pII0002 gene_biotype:protein_coding transcript_biotype:protein_coding description:hypothetical protein
ATGTTTTATGAAGGGAGCAATGCCTCAGCATCAGGTTACGGGGTGACTCACGTAAGGGAC
AGGCAGATGGCAGCTCAGCCACAGGCAGCACTGCAGGAAACTGAATATAAACTGCAGTGA
>CCE57620 cdna plasmid:HUSEC2011CHR1:pHUSEC2011-2:1422:2162:-1 gene:HUS2011_pII0003 gene_biotype:protein_coding transcript_biotype:protein_coding description:site-specific recombinase
ATGAACAATGTCATTCCCCTGCAGAATTCACCAGAACGCGTCTCCCTGTTACCCATTGCG
CCGGGGGTGGATTTTGCAACAGCGCTCTCCCTGAGAAGAATGGCCACTTCCACGGGGGCC
ACACCGGCCTACCTGCTGGCCCCGGAAGTGAGTGCCCTTCTTTTCTATATGCCGGATCAG
CGTCACCATATGCTGTTCGCCACCCTCTGGAATACCGGAATGCGTATTGGCGAAGCCCGG
ATGCTGACACCGGAATCATTTGACCTGGATGGAGTAAGACCGTTTGTGCGGATCCAGTCC
GAAAAAGTGCGTGCGCGACGCGGACGCCCGCCAAAAGATGAAGTGCGCCTGGTTCCGCTG
ACAGATATAAGCTATGTCAGGCAGATGGAAAGCTGGATGATCACCACCCGGCCCCGTCGT
CGTGAACCATTATGGGCCGTGACCGACGAAACCATGCGCAACTGGCTGAAGCAGGCTGTC
AGACGGGCCGAAGCTGACGGGGTACACTTTTCGATTCCGGTAACACCACACACTTTCCGG
CACAGCTATATCATGCACATGCTCTATCACCGCCAGCCCCGGAAAGTCATCCAGGCACTG
GCTGGTCACAGGGATCCACGTTCGATGGAGGTCTATACACGGGTGTTTGCGCTGGATATG
GCTGCCACGCTGGCAGTGCCTTTCACAGGTGACGGACGGGATGCTGCAGAGATCCTGCGT
ACACTGCCTCCCCTGAAGTAA
>CCE57621 cdna plasmid:HUSEC2011CHR1:pHUSEC2011-2:2531:2692:-1 gene:HUS2011_pII0004 gene_biotype:protein_coding transcript_biotype:protein_coding description:plasmid stabilisation system family protein
ATGAGTAATCATAATATCGGTACTCCCCGTCCTGAACTGGGGGAATACACATTCGCACTA
CCCGTTGAACGGCATATGGTTTATTTTCTGCAAACTGATACTGAAATTGTTATTATTCGT
ATATTAAGTCAGCATCAGGATGCCAGCCGTCATTTCAACTGA
>CCE57622 cdna plasmid:HUSEC2011CHR1:pHUSEC2011-2:2894:3076:-1 gene:HUS2011_pII0005 gene_biotype:protein_coding transcript_biotype:protein_coding description:putative transcriptional regulator
ATGGCCAGAACAATGACAGTTGCTCTCGGAGATGAACTCTGGGAGTACATAGAATCTCTC
ATAGAATCAGGTGATTATCGTACACAGAGTGAGGTAATCCGCGAGTCACTTCGTCTCCTT
CGAGAGAAACAGGCAGAGTCACGTCTCCAGGTGTCCTGCGGGATTTACTGGCAGAAGGCT
TGA
>CCE57623 cdna plasmid:HUSEC2011CHR1:pHUSEC2011-2:3145:3420:-1 gene:HUS2011_pII0006 gene_biotype:protein_coding transcript_biotype:protein_coding description:plasmid stabilisation system family protein
ATGGAACTGAAGTGGACCAGTAAGGCGCTTTCTGATTTGTCGCGGTTATATGATTTTCTG
GTGCTGGCCAGTAAACCTGCTGCCGCCAGAACGGTACAGTCCCTGACACAGGCACCGGTC
ATTCTGTTAACTCATCCACGTATGGGAGAACAGTTGTTTCAGTTTGAACCCAGGGAGGTC
AGACGGATTTTTGCTGGCGAGTACGAAATCCGTTACGAAATTAATGGCCAGACTATTTAT
GTATTGCGTCTGTGGCACACACGAGAAAACAGGTAG
>CCE57624 cdna plasmid:HUSEC2011CHR1:pHUSEC2011-2:3420:3698:-1 gene:HUS2011_pII0007 gene_biotype:protein_coding transcript_biotype:protein_coding description:ribbon-helix-helix protein, copG family
ATGAAAAACAATGCCGCACAAGCAACAAAAGTAATTACCGCGCATGTGCCATTACCTATG
GCTGATAAAGTCGACCAGATGGCCGCCAGACTGGAACGTTCCCGGGGCTGGGTTATCAAA
CAGGCGCTTTCTGCATGGCTTGCCCAGGAGGAGGAGCGTAATCGCCTGACGCTGGAAGCC
CTGGACGATGTGACATCCGGACAGGTTATCGACCATCAGGCTGTACAGGCCTGGTCGGAC
AGCCTCAGTACTGACAATCCGTTACCGGTGCCACGCTGA
`

	fastaFile := wtype.File{Name: "ecoli-cdna"}

	err := fastaFile.WriteAll([]byte(fastaText))
	if err != nil {
		t.Fatalf(err.Error())
	}

	database, err := fasta.FastaToDNASequences(fastaFile)
	if err != nil {
		t.Fatalf(err.Error())
	}

	primer := wtype.DNASequence{Seq: "ATGGAACTGAAGTGG"}

	algorithm, found := Algorithms["SWAffine"]
	if !found {
		t.Fatalf("algorithm not found")
	}

	testLimit := 4

	results, err := DNASet(primer, database, algorithm, testLimit)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if len(results) != testLimit {
		t.Errorf("Error, number of results: got %d, want %d\n", len(results), testLimit)
	}

	if results[0].Template.Name() != "CCE57623" {
		t.Errorf("Error, unexpected top scoring result: got %s, want %s\n", results[0].Template.Name(), "CCE57623")
	}

	if results[0].Score() != 15 {
		t.Errorf("Error, unexpected top score: got %d, want %d\n", results[0].Score(), 15)
	}

}

func TestDNASetBenchmark(t *testing.T) {

	t.Skip("requires data file")

	// All ecoli CDNA
	// URI ftp://ftp.ensemblgenomes.org/pub/bacteria/release-42/fasta/bacteria_91_collection/escherichia_coli/cdna/Escherichia_coli.HUSEC2011CHR1.cdna.all.fa.gz
	// ... a few seconds

	dat, err := ioutil.ReadFile("Escherichia_coli.HUSEC2011CHR1.cdna.all.fa")
	if err != nil {
		t.Fatalf(err.Error())
	}

	fastaFile := wtype.File{Name: "ecoli-cdna"}

	err = fastaFile.WriteAll(dat)
	if err != nil {
		t.Fatalf(err.Error())
	}

	database, err := fasta.FastaToDNASequences(fastaFile)
	if err != nil {
		t.Fatalf(err.Error())
	}

	fmt.Printf("Read %d sequences", len(database))

	primer := wtype.DNASequence{Seq: "ATGGAACTGAAGTGG"}

	algorithm, found := Algorithms["SWAffine"]
	if !found {
		t.Fatalf("algorithm not found")
	}

	testLimit := 4

	results, err := DNASet(primer, database, algorithm, testLimit)
	if err != nil {
		t.Fatalf(err.Error())
	}

	for _, result := range results {
		fmt.Printf("Found %s, score %d\n", result.Template.Name(), result.Score())
	}

	if len(results) != testLimit {
		t.Errorf("Error, number of results: got %d, want %d\n", len(results), testLimit)
	}

}
