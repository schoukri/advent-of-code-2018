#!/usr/bin/perl

use strict;
use warnings;

sub top_item_by_value($);

my @lines = <>;
chomp(@lines);


my %asleep;
my %hist;


LINE:
for (my $i = 0; $i < scalar @lines;) {
  my $line = $lines[$i];
  $i++;

  my($id) = $line =~ m/Guard \#(\d+) begins shift/;
  unless ($id) {
    die "can't parse line: $line\n";
  }

  while ($i < scalar @lines) {
    my $sleepLine = $lines[$i];
    my($sleepMin) = $sleepLine =~ m/\:(\d+)\] falls asleep$/;
    next LINE unless defined $sleepMin;

    my $wakeLine = $lines[$i + 1];
    my($wakeMin) = $wakeLine =~ m/\:(\d+)\] wakes up$/;
    unless (defined $wakeMin) {
      die "can't parse wake line: $wakeLine\n";
    }

    $asleep{$id} += ($wakeMin - $sleepMin);
    for (my $min = $sleepMin; $min < $wakeMin; $min++) {
      $hist{$id}{$min}++;
    }

    $i += 2;
  }
}

my($winner_id) = top_item_by_value(\%asleep);

my($top_min) = top_item_by_value($hist{$winner_id});

my $part1 = $winner_id * $top_min;
print "part 1: $part1\n";

my %top_min_per_guard;
my %most_min_per_guard;
while (my($id, $h) = each %hist) {
  my($min, $count) = top_item_by_value($h);
  $top_min_per_guard{$id} = $min;
  $most_min_per_guard{$id} = $count;
}

my($winner2_id) = top_item_by_value(\%most_min_per_guard);

my $part2 = $winner2_id * $top_min_per_guard{$winner2_id};
print "part 2: $part2\n";

sub top_item_by_value($) {
  my $href = $_[0] || die "hashref not specified";
  my($key) = sort {$href->{$b} <=> $href->{$a}} keys %$href;
  return ($key, $href->{$key});
}
