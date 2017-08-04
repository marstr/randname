# Looks at a stream line by line and finds duplicates or incorrect alphabetization.
my $prev = <STDIN>;
while(my $current = <STDIN>) {
    if($current le $prev ){
        exit 1;
    }
    $prev = $current
}
exit 0;