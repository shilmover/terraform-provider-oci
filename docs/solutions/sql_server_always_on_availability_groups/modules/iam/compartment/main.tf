resource "oci_identity_compartment" "compartment" {
  count       = "${var.compartment_id == "" ? 1 : 0}"
  provider    = "oci.home"
  name        = "${var.compartment_name}"
  description = "${var.compartment_description}"
}
