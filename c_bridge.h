#ifndef __C_BRIDGE_H__
#define __C_BRIDGE_H__

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

gpointer
j_async_send(GNetSnmp *session, GNetSnmpPduType type,
             GList *vbl, guint32 arg1, guint32 arg2, GError **error);

gboolean
j_cb_done(GNetSnmp *session, GNetSnmpPdu *spdu, GList *objs, gpointer magic);

void
j_cb_time(GNetSnmp *session, void *magic);

void
j_sync_get(GNetSnmp *snmp, GList *pdu, GError **error);

void
j_sync_send(GNetSnmp *session, GNetSnmpPduType type,
            GList *objs, guint32 arg1, guint32 arg2, GError **error);

void
vbl_delete(GList *list);

#endif //__C_BRIDGE_H__
