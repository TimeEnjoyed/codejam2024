<script lang="ts">

    import Page from "../../components/Page.svelte";
    import {
        Button,
        Card,
        Modal,
        TableBody,
        TableBodyCell,
        TableBodyRow,
        TableHead,
        TableHeadCell, TableSearch, Tooltip
    } from "flowbite-svelte";
    import {onMount} from "svelte";
    import {
        BanUser,
        getAllUsers,
        UnbanUser,
        UpdateDisplayName,
        updateDisplayNameLock
    } from "../../services/adminServices";
    import type {User} from "../../models/user";
    import UserAvatar from "../../components/UserAvatar.svelte";
    import {FontAwesomeIcon} from "@fortawesome/svelte-fontawesome";
    import {faArrowsRotate, faLock, faLockOpen} from "@fortawesome/free-solid-svg-icons";
    import {toast} from "svelte-sonner";

    let users : Array<User> = []; // full list of users
    let displayedUsers : Array<User> = []; // list of Users after being filtered
    let modalOpen : boolean = false;
    let modalUser : User;

    let tableFilter : string = '';
    $: { filterUsers(); }

    function filterUsers() {
        const lowerFilter = tableFilter.toLowerCase();
        displayedUsers = users.filter((item) => {
            return tableFilter.length === 0 ||
                item.DisplayName.toLowerCase().indexOf(lowerFilter) >= 0 ||
                item.ServiceUserId.toLowerCase().indexOf(lowerFilter) >= 0 ||
                item.ServiceUserName.toLowerCase().indexOf(lowerFilter) >= 0 ||
                item.Role.toLowerCase().indexOf(lowerFilter) >= 0 ||
                item.AccountStatus.toLowerCase().indexOf(lowerFilter) >= 0;
        })
    }

    async function loadData() {
        const response = getAllUsers();
        users = await (await response).json() as Array<User>;
        filterUsers();
    }

    function openUserModal(user: User) {
        modalUser = user;
        modalOpen = true;
    }

    function closeUserModal() {
        modalOpen = false;
    }

    function updateUser(user: User) {
        const index : number = users.findIndex((item: User) => { return item.Id === user.Id });
        if (index > -1) {
            users[index] = user;
            filterUsers();
        }
        modalUser = user;
    }

    async function toggleDisplayNameLock() {
        try {
            const updatedUser: User = await updateDisplayNameLock(modalUser.Id, !modalUser.LockDisplayName);
            updateUser(updatedUser);
            toast.success(`User ${updatedUser.DisplayName} Display Name ${updatedUser.LockDisplayName ? 'Locked' : 'Unlocked'}`);
        } catch(err) {
            console.error("ToggleDisplayNameLock:", err);
            toast.error("Error Updating Display Name Lock");
        }
    }

    async function resetDisplayName() {
        try {
            const updatedUser = await UpdateDisplayName(modalUser.Id, modalUser.ServiceUserName);
            updateUser(updatedUser);
            toast.success(`User ${updatedUser.DisplayName} Display Name Reset`);
        } catch (err) {
            console.error("UpdateDisplayName:", err);
            toast.error("Error Resetting Display Name");
        }
    }

    async function banUser() {
        try {
            const updatedUser : User = await BanUser(modalUser.Id);
            updateUser(updatedUser);
            toast.success(`User ${updatedUser.DisplayName} Banned`);
        } catch (err) {
            console.error("BanUser:", err);
            toast.error("Error Banning User");
        }
    }

    async function unbanUser() {
        try {
            const updatedUser = await UnbanUser(modalUser.Id);
            updateUser(updatedUser);
            toast.success(`User ${updatedUser.DisplayName} Unbanned`);
        } catch (err) {
            console.error("UnbanUser: ", err);
            toast.error("Error Unbanning User");
        }
    }

    onMount(() => {
        loadData();
    })

</script>

<Modal bind:open={modalOpen}>
    <svelte:fragment slot="header">
        <UserAvatar class="mr-4" size="md" user={modalUser}/>{modalUser?.DisplayName}
    </svelte:fragment>

    <div class="grid gap-4 grid-cols-[minmax(200px,auto)_1fr]">
        <div class="font-bold">Discord Username</div>
        <div>{modalUser.ServiceUserName}</div>

        <div class="font-bold">Display Name</div>
        <div>
            {modalUser.DisplayName}
            {#if modalUser.Role !== "ADMIN"}
                {#if modalUser.DisplayName !== modalUser.ServiceUserName}
                    <Button class="ml-4 p-2" size="sm" outline on:click={resetDisplayName}><FontAwesomeIcon icon={faArrowsRotate}></FontAwesomeIcon></Button>
                    <Tooltip>Reset Display Name to Discord Username</Tooltip>
                {/if}
                {#if modalUser.LockDisplayName === true}
                    <Button class="p-2" on:click={toggleDisplayNameLock}><FontAwesomeIcon icon={faLock}></FontAwesomeIcon></Button>
                    <Tooltip>User NOT allowed to modify Display Name.</Tooltip>
                {:else}
                    <Button outline class="p-2 text-gray-400" on:click={toggleDisplayNameLock}><FontAwesomeIcon icon={faLockOpen}></FontAwesomeIcon></Button>
                    <Tooltip>User is allowed to modify Display Name.</Tooltip>
                {/if}
            {/if}
        </div>

        <div class="font-bold">Role</div>
        <div>{modalUser.Role}</div>

        <div class="font-bold">Status</div>
        <div>{modalUser.AccountStatus}</div>
    </div>

    <svelte:fragment slot="footer">
        {#if modalUser.Role !== "ADMIN"}
            {#if modalUser.AccountStatus === "ACTIVE" }
                <Button class="border-orange-500 text-orange-500 hover:bg-orange-500" size="sm" outline on:click={banUser}>Ban</Button>
            {:else}
                <Button class="border-orange-500 text-orange-500 hover:bg-orange-500" size="sm" outline on:click={unbanUser}>Unban</Button>
            {/if}
        {/if}
        <div class="w-full flex justify-end gap-4">
            <Button size="sm" outline on:click={closeUserModal}>Close</Button>
        </div>
    </svelte:fragment>
</Modal>

<Page>
    <Card size="xl" class="w-full">
        <h2>Users</h2>
        <div class="flex flex-col gap-8 my-8">
        </div>

        <TableSearch divClass="border-0" innerDivClass="p-4 flex" placeholder="Filter" hoverable={true} bind:inputValue={tableFilter}>
            <svelte:fragment slot="header">
                <span class="relative ml-4 mt-4">{displayedUsers.length} Users</span>
            </svelte:fragment>
            <TableHead>
                <TableHeadCell></TableHeadCell>
                <TableHeadCell>Display Name</TableHeadCell>
                <TableHeadCell>Discord Name</TableHeadCell>
                <TableHeadCell>Discord ID</TableHeadCell>
                <TableHeadCell>Role</TableHeadCell>
                <TableHeadCell>Status</TableHeadCell>
            </TableHead>
            <TableBody>
                {#each displayedUsers as user}
                <TableBodyRow class="border-gray-200 cursor-pointer hover:bg-blue-100" on:click={() => openUserModal(user)}>
                    <TableBodyCell><UserAvatar user={user}/></TableBodyCell>
                    <TableBodyCell>
                        {user.DisplayName}
                        {#if user.LockDisplayName === true}
                            <FontAwesomeIcon icon={faLock}></FontAwesomeIcon>
                        {/if}
                    </TableBodyCell>
                    <TableBodyCell>{user.ServiceUserName}</TableBodyCell>
                    <TableBodyCell>{user.ServiceUserId}</TableBodyCell>
                    <TableBodyCell>{user.Role}</TableBodyCell>
                    <TableBodyCell>{user.AccountStatus}</TableBodyCell>
                </TableBodyRow>
                {/each}
            </TableBody>
        </TableSearch>
    </Card>

</Page>