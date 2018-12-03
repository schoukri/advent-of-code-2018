#!/usr/bin/perl

use strict;
use warnings;

my @data = <>;
chomp(@data);


#1 @ 287,428: 27x20
#2 @ 282,539: 20x10
#3 @ 550,118: 20x23
#4 @ 454,774: 20x19
#5 @ 542,157: 11x24

my $matrix = [];

my %claim;

for my $line (@data) {
  my($id, $x, $y, $w, $h) = $line =~ m/^\#(\d+) \@ (\d+),(\d+): (\d+)x(\d+)$/;

  my $single = [];
  for (my $a = $x; $a < $x+$w; $a++) {
    for (my $b = $y; $b < $y+$h; $b++) {
      $matrix->[$a][$b]++;
      $single->[$a][$b]++;
    }
  }

  $claim{$id} = $single;

}

my $overlap = 0;
for my $row (@$matrix) {
  for my $cell (@$row) {
    next unless defined $cell;
    if ($cell > 1) {
      $overlap++
    }
  }
}

print "part 1: $overlap\n";


# part 2
MATRIX:
while (my($id, $single) = each %claim) {

  for (my $a = 0; $a < scalar(@$single); $a++) {
    my $row = $single->[$a];
    next unless defined $row;
    for (my $b = 0; $b < scalar(@$row); $b++) {
      if (defined $single->[$a][$b] && defined $matrix->[$a][$b]) {
        next MATRIX unless $single->[$a][$b] == $matrix->[$a][$b];
      }
    }
  }

  print "part 2: $id\n";
}

