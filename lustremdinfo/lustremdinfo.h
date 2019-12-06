/*
 * Lustre indexing library for Brindexer
 * Copyright 2019 Cray Inc. All Rights Reserved.
 */

#ifndef _LUSTRE_INDEX_H_
#define _LUSTRE_INDEX_H_

#include <stdint.h>
#include <sys/types.h>
#include <lustre/lustre_user.h>


/* HSM attributes; user-visible copy of Lustre struct hsm_attrs */
struct hsm_user_attrs {
	uint32_t   hsm_compat;
	/* HSM flags, bitfield of Lustre enum hsm_states */
	uint32_t   hsm_flags;
	/* archive id */
	uint64_t   hsm_arch_id;
	/* optional version associated with the last archive operation */
	uint64_t   hsm_arch_ver;
};

/*
 * um_stat	    - stat structure
 * um_valid	    - OBD_MD_FL* flags that show which fields in um_stat are
 *		      valid
 * un_hsm_attrs	    - HSM attributes; hsm_flags and hsm_arch_id might be more
 *                    useful
 * um_numosts       - the number of OSTs the file has objects on. OSTs in
 *                    composite files in which more than one file sublayout has
 *                    objects are counted twice
 * um_mirror_state  - mirror state
 * um_mdt_idx       - MDT index; unimplemnted for now
 * um_ostidx        - the indexes of OSTs the file has objects on; the valid
 *		      values are um_numosts in number. OSTs in composite files
 *		      in which more than one file sublayout has objects are
 *		      counted twice; users of the API can incorporate or ignore
 *		      the duplicate entries
 * um_pool_name     - the pool name for plain layout files that belong to a pool
 *                    the pool name for composite layouts is not returned
 */
struct user_mdinfo {
	lstat_t			um_stat;
	__u64			um_valid;
	struct hsm_user_attrs	um_hsm_attrs;
	size_t			um_numosts;
	enum lov_comp_md_flags	um_mirror_state;
	uint32_t		um_mdtidx;
	uint32_t		um_ostidx[LOV_MAX_STRIPE_COUNT];
	char			um_pool_name[LOV_MAXPOOLNAME + 1];
};

int lustre_get_info_file(const char *path, bool all, struct user_mdinfo *info);

#endif
