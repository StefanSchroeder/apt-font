%has=();
%sho=();
while(<>)
{
	if (m/"(.*)": "(.*)",/)
	{
		#print qq{"$1": "$2",\n};
		$sho{$1}=$2;
		if ($has{$1}) 
		{
			print "\n";
		}
		$has{$1} = 1;
	}
}

foreach my $k (keys %sho)
{
	if ($k =~ m/ /) {
		#print "yes($k)\n";
		my @parts = split(" ", $k);
		pop(@parts);
		my $nk = "@parts";
		print "$nk\n";
		if ($sho{$nk})
		{
			print "Adding $nk for $k\n";
		}
	} else {
		#print "no ($k)\n";
	}
}
