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

// customised g_async_send
gpointer
j_async_send(GNetSnmp *session, GNetSnmpPduType type,
             GList *vbl, guint32 arg1, guint32 arg2, GError **error)
{
    return NULL;
}

// j_cb_done - dummy
gboolean
j_cb_done(GNetSnmp *session, GNetSnmpPdu *spdu, GList *objs, gpointer magic)
{
    CBDone();
    return 1;
}

// j_cb_time - dummy
void
j_cb_time(GNetSnmp *session, void *magic)
{
}

// customised gnet_snmp_sync_get
void
j_sync_get(GNetSnmp *snmp, GList *pdu, GError **error)
{
    gnet_snmp_debug_flags = GNET_SNMP_DEBUG_REQUESTS + GNET_SNMP_DEBUG_SESSION;
    if (gnet_snmp_debug_flags & GNET_SNMP_DEBUG_SESSION) {
        g_printerr("session %p: g_sync_get pdu %p\n", snmp, pdu);
    }
    j_sync_send(snmp, GNET_SNMP_PDU_GET, pdu, 0, 0, error);
}

// customised g_sync_send. No "magic" - use Go concurrency.
void
j_sync_send(GNetSnmp *session, GNetSnmpPduType type,
            GList *objs, guint32 arg1, guint32 arg2, GError **error)
{
    session->done_callback = j_cb_done;
    session->time_callback = j_cb_time;
    if (! j_async_send(session, type, objs, arg1, arg2, error)) {
        if (gnet_snmp_debug_flags & GNET_SNMP_DEBUG_SESSION) {
            g_printerr("session %p: g_sync_send failed to send PDU\n", session);
        }
    }
}

// vbl_delete is wrapper for freeing a var bind list
void
j_vbl_delete(GList *list) {
    g_list_foreach(list, (GFunc) gnet_snmp_varbind_delete, NULL);
    g_list_free(list);
}

