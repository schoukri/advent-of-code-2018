#!/usr/bin/perl

use strict;
use warnings;

use Text::Levenshtein qw(distance);


my @data = <>;
chomp(@data);


for my $one (@data) {
  for my $two (@data) {

    next if $one eq $two;
    next if distance($one, $two) != 1;

    print "one: $one\n"
    print "two: $two\n";

    last;

  }


}
