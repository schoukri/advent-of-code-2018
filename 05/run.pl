#!/usr/bin/perl

use strict;
use warnings;


my @letters = qw(a b c d e f g h i j k l m n o p q r s t u v w x y z);

my @pairs;
for my $l (@letters) {
  push @pairs, uc($l) . $l, $l . uc($l);
}

my $pairs_joined = join '|', @pairs;
my $pairs_regexp = qr/$pairs_joined/;

sub react {
  my $polymer = shift;
  while ($polymer =~ s/$pairs_regexp//go) {}
  return $polymer;
}


my $line = <>;
chomp $line;

print "part 1: " . length(react($line)) ."\n";

my %count;
for my $letter (@letters) {
  (my $copy = $line) =~ s/$letter//gi;
  $count{$letter} = length(react($copy));
}

my($lowest) = sort {$count{$a} <=> $count{$b}} keys %count;

print "part 2: " . $count{$lowest} ."\n";



