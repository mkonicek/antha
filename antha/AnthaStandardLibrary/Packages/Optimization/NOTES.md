#### General

The function `OptimizeAssembly` takes an amino acid translation of a
query sequence (`string`), a multiple sequence alignment (MSA) of
variant DNA sequences (`[]string{}`) and a constraints struct. The
constraints struct contains the DNA query sequence itself (`string`).

The function runs an iterative optimization algorithm (genetic
algorithm) seeking a split of the DNA query (into segments) which:
- satisfies the constraints
- achieves a low value of a certain cost function

`OptimizeAssembly` returns the best split found, with attached overhangs of
the specified length. It would be easy to modify it to return a selection
of close-to-best, alternative splits.

### Mapping to AssemblyProblem
The optimization works with an `AssemblyProblem` and the constraints.

An `AssemblyProblem` consists of an amino acid sequence (`string`) and a
slice of 'mutations' (`[][2]int`).  Each mutation (`[2]int`) consists of
- a position (which must be a codon centre, zero-based) in the underlying DNA
sequence at which mutations will be made
- a number, which is the number of variant codons (distinct) at that position

The mutations data is calculated from the input MSA. If there are no
variants at a particular position, nothing is recorded in the
mutations data. Otherwise, the total number of distinct codons at the
position is recorded.

The code expects (but does not require) that the wild type is
present in the MSA, since this ensures that the wild type codons are
included in the variant counts.

```
Example:

MSA:

agcaagggcgaggagctgttcaccggggtggtgcccatcctggtcgagctggacggcgac - wild type
 S  K  G  E  E  L  F  T  G  V  V  P  I  L  V  E  L  D  G  D
.......t.............c...................................... - variant 1; 2:ggc->gtc, 7:acc->ccc
 .  .  V  .  .  .  .  P  .  .  .  .  .  .  .  .  .  .  .  .
.....................c.t.................................... - variant 2; 7:acc->cct
 .  .  .  .  .  .  .  P  .  .  .  .  .  .  .  .  .  .  .  .

Mutations:
[{7 2} {22 3}]

Mutations without the wild type would be:
[{7 2} {22 2}]

There are 2 variants at position 7, with or without the wild type.
```

### Splits
A split is an `[]int` containing zero-based positions of the split points
in the DNA sequence.

### Constraints
The constraints are:
```
MaxSplits - maximum number of split points 
MinLen - the minimum length of a segment
MaxLen - the maximum length of a segment
MinDistToMut - minimum distance from split point to mutation
EndLen - overhang length
EndsToAvoid - a list of overhang sequences to avoid
Query - the DNA sequence of the query
NoTransitions - prevent possibility of mispairing via transitions
```

### Genetic algorithm (GA)

#### Initialization
A population of candidate splits is generated randomly. Both the
number of split points and the positions of the split points are
chosen randomly.  The size of the population is set by parameter
`pop_size` (default 1000). Splits which violate the constraints are
discarded and the process continues until the required number have
been obtained.

#### Compute cost (fitness)

A cost is calculated for each member (split) of the population.

The cost of a split can be thought of as: the total length of all the
different segments that would be needed, in order to assemble all the
variants.

Example. If there is one mutation site within a segment with two
variants, two different versions of that segment will be needed. If
there is another mutation site within the same segment with three
variants, a total of six different versions of that segment will be
needed - all the combinations of variants at the two sites.

Under the defined cost, splits in which highly mutated sites occur in
the same segment are relatively expensive, because the number of
different versions of the segment that are needed is the product of
the number of variants at the included sites.

SO, the optimization will favour splits which separate mutation sites into
different segments, particularly highly mutated sites close to each
other in the sequence.

Splits in which highly mutated sites occur in long segments are also
relatively expensive, because the cost multiplies the segment length
by the number of variants.

Accordingly, the optimization will favour splits in which highly
mutated sites are in shorter segments.

In general: the cost of a split depends on the positions of the split
points and the distribution of the mutations.

In detail: the cost of a split is calculated as a weighted sum of
segment scores. The unweighted segment score is the length of the
segment.  The weighting multiplier is equal to the product of the
numbers of variants over all mutation sites that fall within the
segment, or 1 if there are no mutation sites within the segment.

```
Example:

Consider assembly problem with mutations
[{7 2} {22 3}]
and DNA query of length 60.

Consider two candidate splits of the interval [0, 60):

> Split 1:
Split [5, 10]
Segment lengths are 6, 6 and 51 (length seems to include ends ???)
Cost is
5 * 1 +       // no mutations, multiplier 1
5 * 2 +       // 1 site with 2 mutations in [5, 10), multiplier 2
50 * 3        // 1 site with 3 mutations in [10, 60), multipler 3
= 157

> Split 2:
Split [5, 40]
Segment lengths are 5, 35 and 20
Cost is
5 * 1 +       // no mutations, multiplier 1
35 * 2 * 3 +  // 1 site with 2 mutations, 1 site with 3 mutations in [5, 40), multiplier 6
20 * 1        // no mutations, multiplier 1
= 235

```

#### Evolution loop
The population of splits is evolved 1000 times (`max_iterations`) to obtain
the best split. In each iteration population is regenerated via a process of
selection, recombination and mutation events. Low scoring splits are
preferentially 'selected' for the new population. There is a
probability for splits to 'recombine' (a new split is constructed by
sampling points from two parent splits) or 'mutate' (split points are
added, removed or shifted by some amount). Throughout all events,
constraints are checked and any candidate splits violating the
constraints are removed.

### Contexts other than protein mutagenesis
At first sight the existing code looks tightly coupled to the context of
protein mutagenesis. However (with a couple of minor edits) it runs fine
with just one sequence in the input MSA. In this case the AssemblyProblem ends
up with an empty mutations slice, and the cost takes a uniform value independent of
split: namely, the total length of the DNA sequence. It appears that none
of the constraints checking relating to protein mutagenesis will be
triggered when there are no mutations. But other constraints, such as the
max/min segment length and max splits, will continue to operate. So overall
the code looks pretty flexible and adaptable to other contexts.

### A ligation fidelity model
Currently, the cost function takes no account of differences in the
ligation fidelity of overhangs. It's possible to avoid certain
overhangs via the constraints parameter, but apart from this, all
allowed overhangs are treated equally.

A ligation fidelity model (LFM) would predict: the probability `p` of a correct
assembly, given the set of ends implied by a candidate split. It would be a 
probabilistic model based on data on rates of correct pairing and mis-pairing 
for the different pairs of ends.

It appears straightforward to pass required data into the cost function: it would be added
to the existing  `AssemblyProblem`. This would a simple and low-risk change,
easy to switch off in order to recover current behaviour.

The LFM probability could modify the existing cost in the following way:
```
cost -> cost * ( 2 - p )
```

If there are no mutations, cost is constant (query length) so the
optimization will choose a split which minimizes the probability
of mis-assembly.

Where there are mutations, the existing cost is multiplied the factor `2 - p`
which favours splits with lower probability of mis-assembly, but remains
positive even if `p = 1` allowing the cost to exert an effect 
in all cases. The formula reflects a particular trade-off between getting a low 
(unmodified) cost and a low probability of mis-assembly. Other trade-offs 
could be used.