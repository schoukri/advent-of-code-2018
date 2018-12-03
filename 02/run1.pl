#!/usr/bin/perl

use strict;
use warnings;

my @data = <>;
chomp(@data);

my $two = 0;
my $three = 0;

for my $line (@data) {

  my %count;
  for my $letter (split '', $line) {
    $count{$letter}++;
  }

  if (grep {$_ == 2} values %count) {
    $two++;
  }

  if (grep {$_ == 3} values %count) {
    $three++;
  }

}

printf "answer: %d\n", $two * $three;
