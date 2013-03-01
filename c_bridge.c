#include "c_bridge.h"

#include <gsnmp/ber.h>
#include <gsnmp/pdu.h>
#include <gsnmp/dispatch.h>
#include <gsnmp/message.h>
#include <gsnmp/security.h>
#include <gsnmp/session.h>
#include <gsnmp/transport.h>
#include <gsnmp/utils.h>
#include <gsnmp/gsnmp.h>
#include <stdlib.h>

// get_err_label is a wrapper for gnet_snmp_enum_get_label()
gchar const *
get_err_label(gint32 const id) {
	return gnet_snmp_enum_get_label(gnet_snmp_enum_error_table, id);
}

// vbl_delete is a wrapper for freeing a var bind list
void
vbl_delete(GList *list) {
	g_list_foreach(list, (GFunc) gnet_snmp_varbind_delete, NULL);
	g_list_free(list);
}
