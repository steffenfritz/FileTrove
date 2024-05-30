rule TestTextRule
{
    strings:
        $my_text_string = "worm"

    condition:
        $my_text_string
}

rule TestPNGRule
{
    strings:
      $my_hex_string = { 50 4e 47 }
    
    condition:
      $my_hex_string

}
